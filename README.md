GrpcAPI
========
by beji
--------------

# Install

<pre><code> go get github.com/jooita/grpcapi </code></pre>


# Usage

* Yaml Config File
Default file name is 'grpcapi.yaml'.
This is a example.
<pre><code>  
api-server:
        dockerdir: ./api-server
        env: env/api-server.env
        apiport: 50051
        apidir: /api

api-client:
        dockerdir: ./api-client
        env: env/api-client.env
        lang:
         - go
         - python
        command: ls -l /api
</code></pre>

* Execute
<pre><code> cd $GOPATH/src/github/jooita/grpcapi && go build </code></pre>

<pre><code> ./grpcapi up </code></pre>
