# Tracer

```javascript
"entity": {
  "id": "1234", // required, can be the same value as uri field
  "uri": "http://new_doc.com", // can be ommited if same value as id
  "name": "Some Document", // optional, recommended for readability
  "creationDate": "2016-06-19", // optional, recommended for readability
  "type": "document", // optional, recommended for readability
  "data": {}, // optional, collection of additional attributes and values
  "wasGeneratedBy": {}, // required, activity that generated this entity
  "wasDerivedFrom": [] // list of entities the generated entity derives from, can be omitted if same as uesd field in generating activity
}
```

```javascript
"activity": {
  "id": "5678", // required
  "type": "batch", // required
  "isBatch": true, // required
  "name": "Some Activity", // optional, recommended for readability
  "startDate": "2016-06-18", // optional, recommended for readability
  "endDate": "2016-06-19", // optional, recommended for readability
  "data": {}, // optional, collection of additional attributes and values
  "wasAssociatedWith": {}, // required, agent responsible for activity
  "used": [] // list of entities used in the activity, can be omitted if activity type is batch.
}
```

```javascript
"agent": {
  "id": "9876", //required
  "name": "Exhauster", // optional, recommended for readability
  "type": "service", // optional, recommended for readability
  "data": {}, // optional, collection of additional attributes and values
  "actedOnBehalfOf": {} // optional, supervisor responsible for acting agent, can be omitted if unsupervised
}
```

```json
{
  "entity": {
    "id": "55",
    "uri": "http://new_doc.com",
    "name": "Some Document",
    "creationDate": "2016-06-19",
    "type": "document",
    "wasDerivedFrom": [
      {
        "id": "2"
      }
    ],
    "data": {
      "revision": 2,
      "description": "...",
      "author": "Max Mustermann"
    },
    "wasGeneratedBy": {
      "id": "blablabla",
      "name": "Some Activity",
      "startDate": "2016-06-18",
      "endDate": "2016-06-19",
      "type": "aggregation",
      "isBatch": true,
      "data": {
        "errors": 0,
        "duration": "12h"
      },
      "wasAssociatedWith": {
        "id": "9876",
        "name": "Exhauster",
        "type": "service",
        "data": {},
        "actedOnBehalfOf": {
          "id": "5432",
          "name": "Hans",
          "type": "contractor",
          "data": {
            "adress": "...",
            "email": "...",
            "phone": "..."
          }
        }
      },
      "used": [
        {
          "id": "1"
        }
      ]
    }
  }
}
```
