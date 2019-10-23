[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2FAriel-dono%2FInvoice_Test.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2FAriel-dono%2FInvoice_Test?ref=badge_shield)

### Sobre la solución implementada:
##### Framework:
- Gorilla Mux
##### Base de datos:
- Couchbase
##### Dependencias externas:
- "gopkg.in/couchbase/gocb.v1"
- "github.com/google/uuid"
- "github.com/gorilla/mux"
##### Mejoras:
- La carga/llamada al servicio del BCCR debería realizarse de forma aislada una
  vez por día por lo que optaría por crear un servicio independiente con un job
  que interactue con esta solución en cuestión.
- Implementar configuración parametrizable

###### Ampliar en consideraciones sobre couchbase:
https://docs.couchbase.com/go-sdk/1.6/start-using-sdk.html

###### Ampliar sobre N1QL y Mutations:
https://docs.couchbase.com/go-sdk/current/n1ql-queries-with-sdk.html
https://docs.couchbase.com/go-sdk/current/subdocument-operations.html

###### Ampliar en consideraciones sobre Gorrilla Mux:
http://www.gorillatoolkit.org/pkg/mux

###### Ampliar sobre la implementación de UUID de google:
https://godoc.org/github.com/google/uuid
https://github.com/google/uuid


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2FAriel-dono%2FInvoice_Test.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2FAriel-dono%2FInvoice_Test?ref=badge_large)