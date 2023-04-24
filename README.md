# Microsserviço de conversão de vídeos

No projeto _Codeflix_, o microsserviço de conversão de vídeos é uma aplicação que tem, como objetivo, converter vídeos a partir do formato _MP4_ para o formato _MPEG-DASH_, o qual corresponde a um formato mais adequado para realizar o _playback_ de vídeos na _Internet_.

> Por que Go?

- Optou-se pela linguagem _Go_ porque ela lida de maneira muito simples com a complexidade técnica que a aplicação precisa resolver:
  - O _Go_ trabalha muito facilmente com _multithreading_;
  - É muito performático;
  - E vai realizar a leitura de filas, o _upload_ dos vídeos e o _encoding_ de mais baixo nível de uma forma muito mais simples.

O fluxo da aplicação, de forma geral, consiste em:

1. Receber uma mensagem via _RabbitMQ_;
2. Fazer o _download_ do vídeo no _Google Cloud Storage_;
3. Fragmentar o vídeo;
4. Converter o vídeo para _MPEG-DASH_;
5. Fazer o _upload_ do vídeo para o _Google Cloud Storage_;
6. Enviar uma notificação via fila com informações do vídeo convertido ou informando erro na conversão;
7. Em caso de erro, a mensagem original enviada via _RabbitMQ_ é rejeitada e encaminhada diretamente a uma _Dead Letter Exchange_.

Um cenário comum de erro pode ser quando, por exemplo: 1. Um formato incorreto de dados for enviado para a fila de entrada; 2. O vídeo informado para download não existir; 3. Ocorrer um erro no processo de conversão de um vídeo.

Então, o que acontece? O microsserviço envia uma mensagem de notificação contendo a mensagem de erro.

E quanto à mensagem originalmente enviada pelo _RabbitMQ_?

Para não perder a mensagem, é dado um _reject_ nela. Assim, toda mensagem rejeitada na fila de entrada é encaminhada para uma _Dead Letter Exchange_. A _Dead Letter Exchange_ encaminha, automaticamente, para uma _Dead Letter Queue_, que fica aguardando todas as mensagens que tiveram problema. Pode ser feito, então, uma consulta nessa fila para uma análise manual, procurando entender por que ocorreu o problema.

#### Como funciona o processamento?

Durante o processamento, a aplicação processa diversas mensagens de forma paralela/concorrente. E um simples arquivo de _MP4_, quando convertido para _MPEG-DASH_, é segmentado em múltiplos arquivos de áudio e vídeo. Logo, o processo de _upload_ não ocorre para apenas um único arquivo, porque ele também acontece de maneira paralela/concorrente.

Mais especificamente, quando é consumida uma nova mensagem, o que acontece? O microsserviço cria uma nova _thread_, a qual o _Go_ chama de _goroutine_. Essa thread é responsável por processar o vídeo e, depois de convertê-lo, é feito o _upload_. Assim, essa thread atua como um _job manager_ ou _worker_, gerando mais n _threads_ que vão realizar o _upload_ para subir diversos arquivos de uma só vez.

Assim, quando a aplicação estiver lidando com o processamento de vários vídeos, é possível parametrizar o número máximo de _threads_ que vão iniciar o processo de conversão dos vídeos. Também é possível parametrizar o número máximo de _threads_ que vão fazer o _upload_ dos múltiplos fragmentos gerados quando o vídeo já estiver convertido.

#### E o design da aplicação?

Optou-se por aproximar-se o máximo possível de uma Arquitetura Hexagonal (_Ports And Adapters_).

> Arquitetura Hexagonal (_Ports and Adapters_)

  - #### Permite trabalhar com um _design_ focado em solucionar o problema do domínio;

    - Vamos ter uma camada de domínio responsável por resolver a complexidade do negócio;
  
  - #### Permite deixar a complexidade técnica para ser resolvida por uma camada de _framework_;

    - Vamos ter uma camada de _framework_ responsável por resolver o sistema de mensageria (_RabbitMQ_) e banco de dados (_Postgres_);
    
- Com isso:

    - A  aplicação torna-se flexível para adicionar/remover componentes de infraestrutura sem precisar alterar _nenhum_ outro componente da aplicação ou o modelo de domínio;

Então, a aplicação divide-se, basicamente, em 3 camadas: _Domain_, _Application_ e _Framework_.

_Domain_ corresponde ao coração da aplicação, sendo composto por entidades e regras de negócio. _Application_ vai corresponder aos casos de uso, onde se utiliza o _Domain_ para executar o fluxo da aplicação.

E a última camada, chamada de _Framework_, corresponde ao conjunto de bibliotecas que vão dar acesso à aplicação. Compõe-se, por exemplo, de bibliotecas que vão possibilitar receber mensagens das filas, conectar com o banco de dados, etc.

### Execução

1. Executar: `docker-compose up -d`;

#### RabbitMQ

2. Acessar: `localhost:15672`;

3. Criar uma fila para as notificações de sucesso: _videos-result_ e vinculá-la ao _Exchange_: _amq.direct_ / _Routing Key_: _jobs_;

4. Criar uma nova _Exchange_: _dlx_ / _Type_: _fanout_;

5. Criar uma nova fila para as notificações de erro: _videos-failed_ e vinculá-la ao _Exchange_: _dlx_;

#### Golang

6. Entrar no _container_ da aplicação: `docker-compose exec app bash`;

7. Subir a aplicação: `go run framework/cmd/server/server.go`;

#### RabbitMQ

8. Acessar a fila _videos_ / _Publish message_:

```
{
  "resource_id":"id-client1",
  "file_path":"convite.mp4"
}
```

9. Onde _file_path_ equivale ao arquivo armazenado no _bucket_ do _Cloud Storage_ para conversão;

10. Clicar _Publish message_ 5 vezes publica 5 mensagens para serem consumidas pela aplicação _Golang_. A aplicação é parametrizada com 5 _workers_ (_threads_) para o processo de conversão e cada _worker_ conta com 50 _threads_ para o _upload_ dos fragmentos já convertidos. Ao todo, 1250 _threads_ operam de forma paralela/concorrente;

#### Golang

11. Visualizar o processamento no terminal. Como o banco de dados _Postgres_ está em modo _debug_, é possível visualizar o _status_ de cada _job_ sendo alterado para: _DOWNLOADING_, _FRAGMENTING_, _ENCODING_, _UPLOADING_, _FINISHING_, _COMPLETED_;

#### RabbitMQ

12. Acessar a fila _videos-result_ para visualizar as notificações de sucesso;

13. Acessar _Get message_ / clicar _Get Message(s)_;

#### Google Cloud Storage

14. Visualizar 5 novos diretórios, cada diretório contendo o resultado da conversão para o formato _MPEG-DASH_.
