## Documentation

### How score is calculated

* **For a single category:**

   - **Weighted Average:** 

$$\text{Weighted\_Average\_Category} = \frac{\sum (Rating_i \cdot Weight_i)}{\sum Weight_i}$$

      * where:
         * `Rating_i`: Rating for the i-th review within the category.
         * `Weight_i`: Weight associated with the category of the i-th review.

* **For all categories:**

   - **Overall Weighted Average:**
$$\text{Overall\_Weighted\_Average} = \frac{\sum (\text{Weighted\_Average\_Category}_j)}{\text{Number\_of\_Categories}}$$

      * where:
         * `Weighted_Average_Category_j`: Weighted average for the j-th category.
         * `Number_of_Categories`: Total number of rating categories.

**2. Score Calculation:**

* **Normalized Score:** 

$$\text{Score} = \left(\frac{\text{Overall\_Weighted\_Average}}{\text{Maximum\_Possible\_Rating}}\right) \times 100$$

   * where:
      * `Maximum_Possible_Rating`: The highest possible rating value ir our case 5

### Testing Locally

For testing server locally, you can use docker-compose file:

    1. `docker-compose -f ./docker-compose.yaml  build`
    2. `docker-compose -f ./docker-compose.yaml  up`

Project includes tests for score service and for grpc server in test folder.

### Suggested way of deployment

Regarding deployment, I would use a Helm Charts, so we have centralize the infrastructure of the service in the source code, we can version it, and publish as we would do for the service.

A possible flow should look like this:
   
1. Developer commit into master branch.
2. CI pipeline e.g Jenkins, Github actions or Bitbucket Pipeline, builds docker image and helm charts for the service, both are tagged with semantic version and push to company docker registry.
3. A GitOps tool like Flux or Argo CD monitors a central Git repository containing the desired state of the Kubernetes cluster.

    When a new version of the Helm Chart is detected in the repository, the GitOps tool automatically updates the Kubernetes cluster to reflect the desired state, deploying the new application version.

I have included a minimal configuration for Helm Chart, it can be use locally with minikube, following these steps:

1. `docker-compose -f ./docker-compose.yaml  build` -> To build docker image, this will create a docker image with the following tag score-service:latest
2.  If your minikube can see images from your local docker you will be able to use the image if not, then you need to upload image to minikube using minikube cache command.
3.  Install score-service using helm cli, with the following command `helm install score-service --namespace score-service --create-namespace ./helm`
4.  You should be able to see a new namespace created in your cluster called score-service with a pod for the service.

NB! Currently the service gets deployed successfully, however the database is not present. The reason is that mounting a SQLite database is a bit more complicate to do and there are different approaches.



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
