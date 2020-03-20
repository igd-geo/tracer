# Tracer
[![Apache License, Version 2.0](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0)

**Disclaimer. This project is a work in progress, hence errors and breakage is to be expected and 
therefore not suitable for production.**

Tracer is a provenance management service built to passively generate provenance by listening to 
dedicated messages sent by services during the execution of scenarios within the marauder ecosystem.

## Features
Tracer has the following features:
* Collect information required to generate provenance directly from RabbitMQ exchanges
(requires messages to be sent with a spefific routing key)
* Parse collected information into a provenance data model derived from W3C PROV-DM
* Batching of new database entries to improve write performance
* Efficient storage and retrieval of provenance by combining the data models graph-like structure
with Dgraph, a modern graph database
* Queryable API through GraphQL

## Roadmap
The following features are planned for tracer:
* Implement a [zap](https://github.com/uber-go/zap) logger to improve logging, e.g. send log 
messages to stdout and RabbitMQ simultaneously
* Improve transaction handling to improve write performance
* Migrate to Dgraph 2.0 to make use of its new native GraphQL features to increase the GraphQL 
API read performance
* Implement connection pools to collect provenance from multiple exchanges at the same time
* Support multiple methods to collect provenance for services that may not provide dedicated 
provenance messages:
  * Collect provenance by listening and filtering messages send to the system log
  * Retroactively collect provenance information by querying registried provided by the mrauder 
  ecosystem.
* Add support for other database-backends:
  * neo4j

# Installation
## From Source
The binaries required for the client that generates and stores provenance and the binary for the
API webserver can be installed using `go get`. Private repositories require some
[adjustments](https://golang.org/doc/faq#git_https) to be able to access the packages.<br>
**Note.** This installation method requires the `config.yaml` files, see 
[Configuration](#Configuration).
```
go get geocode.igd.fraunhofer.de/hummer/tracer/cmd/client
go get geocode.igd.fraunhofer.de/hummer/tracer/cmd/api
```

Alternatively the project can be build by cloning the repository and using `go build` or
`go install`, e.g.:
```
go build -o tracer-client ./cmd/client
go build -o tracer-client ./cmd/api
```

## Using Docker
The project provides prebuilt docker images for both the client and the API. They can be used as
follows:
```
docker run --rm geocode.igd.fraunhofer.de:4567/hummer/tracer/client
docker run --rm geocode.igd.fraunhofer.de:4567/hummer/tracer/api
```

## Using Docker Compose
Lastly the project provides a full docker compose configuration that includes Tracer and all its
dependencies:
```
docker-compose up
```

# Usage
## Configuration
The application can be configured using the `config.yaml` files in `/cmd/client` and `/cmd/api`.
The configuration files need to be in the same directory as the binary. When using the docker
images, the configuration files can be mounted at runtime:
```
docker run -v /path/to/config.yaml:/tracer/config.yaml --rm geocode.igd.fraunhofer.de:4567/hummer/tracer/client
docker run -v /path/to/config.yaml:/tracer/config.yaml --rm geocode.igd.fraunhofer.de:4567/hummer/tracer/api
```
## Generating Provenance
A Service that wishes for provenance to be generated needs to send dedicated provenance messages.
Tracer assumes that the messages are being sent by processes that represent a service within
the marauder ecosystem. Therefore provenance messages need to be routed through RabbitMQ with a
specific routing key.<br>
**Note.** This assumes that the exchange used to route the messages uses the **Scenario ID** as 
name.

They routing key has to be structured as follows:
```
<Service ID>.<Process ID>.provenance
```

Provenance messages must have a JSON formatted body.

**Note.** Currently Tracer can only generate provenacne for processes with a single output, however,
multiple inputs are supported.
```
{
  "timestamp": "2000.01.01 12:12:12.000",
  "input": "["entityUID"]",
  "output": "entityUID"
}
```

## Provenance Components
Tracer uses a data model derived from W3C PROV-DM, see [W3C PROV](https://www.w3.org/TR/prov-dm/)
for more, and therefore divides provenance into three main components:

**Entities.** Entities have the following queryable attributes:
  * uid
  * id
  * uri
  * type
  * name
  * creationDate

and the following edges:
   * wasDerivedFrom (points to zero or more entities)
   * wasGeneratedBy (points to an activity)

**Activities.** Activities have the following queryable attributes:
  * uid
  * id
  * type
  * name
  * startDate
  * endDate

and the following edges:
 * wasAssociatedWith (points to an agent)
 * used (points to zero or more entities)

**Agents.** Agents have the following queryable attributes:
  * uid
  * id
  * type
  * name

and the following edges:
  * actedOnBehalfOf (points to an agent)

**Note.** More attributes will be added as needed.

## GraphQL
Tracer offers a queryable API using GraphQL.
Provenance components can be retrieved using the GraphQL endpoint, e.g. with HTTPie:
```
http 'localhost/graphql?query={entity(id:""){uid,id,uri,type,name,creationDate}}'
http 'localhost/graphql?query={activity(id:""){uid,id,type,name,startDate,endDate}}'
http 'localhost/graphql?query={agent(id:""){uid,id,type,name}}'
```

In addition to simple components, GraphQL is able to traverse the provenance graph using the edges
as described above, e.g. the following query returns an entity, the activity that generated it and
the agent that is responsible for that activity:
```
http 'localhost/graphql?query={entity(id:""){uid,id,uri,type,name,creationDate,wasGeneratedBy{uid,id,type,name,startDate,endDate,wasAssociatedWith{uid,id,type,name}}}}
```
Formatted:
```
{
  entity(id:"") {
    uid,
    id,
    uri,
    type, 
    name, 
    creationDate
    wasGeneratedBy{
      uid,
      id,
      type,
      name,
      startDate,
      endDate,
      wasAssociatedWith{
        uid,
        id,
        type,
        name
      }
    }
  }
}
```

# License

Tracer is licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for the full license text.

Icons made by [Freepik](https://www.flaticon.com/authors/freepik "Freepik") from [www.flaticon.com](https://www.flaticon.com/ "Flaticon")
