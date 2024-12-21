Documentation:

- Ticket score algorithm:
    Calculates the weighted average rating for each ticket by:
    -   Multiplying each rating by its corresponding category weight.
    -   Summing these weighted ratings for each ticket.
    -   Dividing the sum of weighted ratings by the sum of category weights for that ticket.
    With the weighted average then:
    -   Divides the weighted_average for each ticket by 5 (maximum possible rating).
    -   Multiplies the result by 100 to convert it to a percentage.

    This is used a based calculation to all queries used to generate the results for each endpoint, queries can be found in repository folder.
- For testing server locally, you can use docker-compose file:
    1. `docker-compose -f ./docker-compose.yaml  build`
    2. `docker-compose -f ./docker-compose.yaml  up`

- Project includes tests for score service and for grpc server in test folder.

- Regarding deployment, my suggested approach is to use a Helm Charts, so we have centralize the infrastructure of the service in the source code. The flow should be like this:
    1. Developer push commit into master branch
    2. CI pipeline e.g Jenkins, Github actions or Bitbucket Pipeline, should build docker image, tag it with the correct version and push the image to a docker registry amd also Helm Chart for service could be deploy to same registry.
    3. Using a GitOps tool like Flux CD, a central repository with relevant configuration for Helm Chart, could detect new version of Helm Chart of service to detect if any updates are needed for infra or/and new version of application was deployed, and deploy changes if needed to K8S cluster.


# Software Engineer Test Task

As a test task for [Klaus](https://www.klausapp.com) software engineering position we ask our candidates to build a small [gRPC](https://grpc.io) service using language of their choice. Preferred language for new services in Klaus is [Go](https://golang.org).

The service should be using provided sample data from SQLite database (`database.db`).

Please fork this repository and share the link to your solution with us.

### Tasks

1. Come up with ticket score algorithm that accounts for rating category weights (available in `rating_categories` table). Ratings are given in a scale of 0 to 5. Score should be representable in percentages from 0 to 100. 

2. Build a service that can be queried using [gRPC](https://grpc.io/docs/tutorials/basic/go/) calls and can answer following questions:

    * **Aggregated category scores over a period of time**
    
        E.g. what have the daily ticket scores been for a past week or what were the scores between 1st and 31st of January.

        For periods longer than one month weekly aggregates should be returned instead of daily values.

        From the response the following UI representation should be possible:

        | Category | Ratings | Date 1 | Date 2 | ... | Score |
        |----|----|----|----|----|----|
        | Tone | 1 | 30% | N/A | N/A | X% |
        | Grammar | 2 | N/A | 90% | 100% | X% |
        | Random | 6 | 12% | 10% | 10% | X% |

    * **Scores by ticket**

        Aggregate scores for categories within defined period by ticket.

        E.g. what aggregate category scores tickets have within defined rating time range have.

        | Ticket ID | Category 1 | Category 2 |
        |----|----|----|
        | 1   |  100%  |  30%  |
        | 2   |  30%  |  80%  |

    * **Overall quality score**

        What is the overall aggregate score for a period.

        E.g. the overall score over past week has been 96%.

    * **Period over Period score change**

        What has been the change from selected period over previous period.

        E.g. current week vs. previous week or December vs. January change in percentages.


### Bonus

* How would you build and deploy the solution?

    At Klaus we make heavy use of containers and [Kubernetes](https://kubernetes.io).
