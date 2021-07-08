### 12 Factors and Cloud-Native Microservices

#### Looking at Cloud-Native services through the 12 Factors’ lens 🔎

![](https://cdn-images-1.medium.com/max/1200/1*Jxs4kBmyccdCiIkxBmKIkA.png)
<span class="figcaption_hack">12 Factors to be considered for Cloud-Native Microservices as Kubernetes
Containers</span>

**W**e are in the world of [Everything As A
Service](https://simple.wikipedia.org/wiki/Everything_as_a_service)(*EaaS*)
which is becoming more and more common now. When the software is provided as a
service, it implies that a few features viz., security, scalability,
reliability, and performance must be implicit.

To meet this, the cloud-native applications are developed. A cloud-native app is
best defined as the applications that has set of characteristics which
differentiates from traditional applications. The cloud-native applications are
*microservices* in architectural nature. They comply with the principles of *12
Factor App*. Apps that are developed with the 12 Factor principle *can run in
any PaaS* infrastructure. All the service-to-service communicatio happen via
*APIs*. These apps can *scale on demand*.

![](https://cdn-images-1.medium.com/max/800/1*4eo06nidlNCTVejtnBAeSQ.png)
<span class="figcaption_hack">12 Factor site banner; Source: 12factor.net</span>

In this article we will see how to apply the set of [twelve
factors](https://12factor.net/)** **while designing and developing a modern,
cloud-native application.

### *I. Codebase*

> One codebase tracked in revision control, many deploys.

The 12-Factor app advocates that every application or service should have its
own codebase in a version control system like git. Even when building for
cross-platform like Mobile, Windows, Linux, Mac, or any others, there should not
be multiple codebases for the same application or service. This principle must
be followed for the deployment of the application or service in multiple
environments as well viz., dev, staging, production, etc.

When there’s a need for sharing the same codebase for multiple applications,
then that shared code should be modeled as a library and the library should have
its own repository. This makes sure that there is a one-one correlation between
codebase and application or service.

Multiple branches should be created to work on various features, fixes,
enhancements of the application or service in the version control system.
Storing the codebase in a version control system paves an elegant way to resolve
the merge conflicts, and also the ability to roll back the code to previous
historical versions. It's a place from where teams can plan to do continuous
integration (CI) and continuous deployment (CD).

![](https://cdn-images-1.medium.com/max/800/1*3l_QarLmVNq_zh2GG7pJTQ.png)
<span class="figcaption_hack">Multiple git branches and one repository with CI/CD to multiple environments</span>

### II. Dependencies

> Explicitly declare and isolate dependencies.

These are two considerations on dependencies with the 12 Factor App. Declaring
the dependencies and isolating them.

By doing so, one of the main problems can be eliminated viz., *It works on my
machine* scenario. Adding to explicit declaration of the dependencies, adding
the dependency version locking files if available for the language like
`go.sum`, `package-lock.json`, `yarn.lock`, etc., will make sure to stick with
the same version of the dependencies while building your application anywhere.
This will also make sure that all the environments are in sync and reproduce the
same behavior.

Since the 12 Factor App recommendation is one codebase per application, its’ the
teams’ responsibility to make sure no one is merging any upgraded version of the
dependencies accidentally to the main branch without testing.

On the Isolation part, the dependencies should be always pulled from the
language corresponding registries. These must not be packaged as part of the
application. The 12 Factor App also considers the usage of OS utilities as
dependencies too. For instance, if there is a need for using the `sleep, wget,
`etc., are also considered as dependencies. Even though these may present in the
OS, the application developer must make sure everything is available for the
application to fulfill the functionality.

For a containerized environment, the Dockerfile is used for adding the os level
dependencies as well.

![](https://cdn-images-1.medium.com/max/800/1*DtDDeBqLDGnC0ERH3EyB-A.png)
<span class="figcaption_hack">screenshot of package-lock.json and go.sum file</span>

### III. Configuration

> Store the config in Environment Variable.

All modern app needs some configurations. These must be always separated from
the functional code. Anything that varies from environment to environment such
as development, test, and production is considered as configuration. There
should be a strict separation between config and code, as these are usually
service account credentials, resource handles to backend databases or the
application runtime properties. The code should remain the same irrespective of
where the application is being deployed, but these configurations can vary.

Also, having an isolated the configuration from the code simplifies the process
of updating config values. This also helps in eliminating the need for
re-deployment of the application when certain config values are changed.

In the *Kubernetes* paradigm, we have various options for separating
configurations from the code. One can use `configMap` which holds the key-value
pair in a plain text format. A sample configMap will look like the below.

In the above sample configMap, we have three different sets of properties viz.,
`app_mode`, `game.properties` and `ui.properties`. This can be used as either
environment variables or files inside the container and make the code refer to
these.

In the case of an environment variable, the container’s spec will include the
snippet below

With the above configuration of env in the deployment of a pod, the env variable
`APP_MODE`in the container will be equated to the value of the variable
`app_mode` in the config map as shown below.

![](https://cdn-images-1.medium.com/max/800/1*JaN726ZUbNRpE2CCgqfp7A.png)
<span class="figcaption_hack">Configuration value injected into the container</span>

Environment Variables might be the right fit for all the use-cases. In case of
accessing the *configMap* data as a file from within the container, you mount
that as data volumes in the deployment as shown below.

With the above snippet, all the keys in the config map will be referred to as a
file from within the container as shown below.

![](https://cdn-images-1.medium.com/max/800/1*rRg-1sxMQy_sIQgXy4nLGQ.png)
<span class="figcaption_hack">configMap values mounted to the container as files</span>

These are not the only way to access the configurations. You can also access
these configurations in a more cloud-native way as well. For example, if you
want to access a configuration stored in Kubernetes configMap via a *GoLang*
microservice, you can include this
[k8s-discovery](https://pkg.go.dev/github.com/gkarthiks/k8s-discovery) module
and pass the configMap. More details on this can be found
[here](https://medium.com/swlh/clientset-module-for-in-cluster-and-out-cluster-3f0d80af79ed).
In this way, we don’t need to wait for the changes if there are any that
occurred to the values in the configMap to reflect into the container as this
module queries the Kubernetes API for every reference.

These are the configurations for the applications that can be stored in the git
repo as plain text even. But, to have some sensitive data separated from the
functional code, Kubernetes provides a way to encode and store it as binary
data. These are called *secrets*. These secrets can also be mounted as
environment variables and as a file in the same way as a configMap. For example,
you have a username and a password as shown below.


Now you want to add it to the environment variable of the microservice
container. For that, we need to create the secret in Kubernetes first. Refer to
the [k8s docs
here](https://kubernetes.io/docs/tasks/configmap-secret/managing-secret-using-kubectl/)
for ways to manage secrets. After creating the secret in the Kubernetes, the K8s
resource YAML will look like below.

Our username and password are encoded with base64 and stored there. To get the
plain data out of that execute the below command in the terminal and verify.


Now that we have the secret, let's inject it into the pod and verify it. Adding
secret as an environment variable and volume will look like the below.

Now if you exec into the container and see, you can see the decoded value in the
files and environment variables as shown below.

![](https://cdn-images-1.medium.com/max/800/1*-Pd-txmxhxE8DRI0sRafXA.png)
<span class="figcaption_hack">secrets are injected and mounted to the container as env values and files</span>

Now there is a problem. Yes, we don’t want sensitive information to be exposed
like this. Anyone having the access to the container can see this data. As this
is only *base64* encoded and also decoded while injected into the container. To
avoid this, the best way to store sensitive data like passwords, certificates,
keys, etc., is to use a secure vault outside of the container. Many such
solutions are available like *Hashicorp Vault*, *CyberArk Conjur*, *AWS Secrets
Manager*, *Azure Key Vault, to* name a few. And access to these services should
be provided on a token basis which has a TTL attached to it.

### IV. Backing Services

> Treat backing services as attached resources.

Anything required to be accessible over the network to fulfill the application's
functionality is considered backing services. These are mostly services like
databases, secret stores, messaging queues, caching services, etc., making sure
to isolate these backing services segregated from the core functional code and
treating them as an attached service will help in easy portability. This also
helps in swapping the resources from one service provider to another,
eliminating the vendor-lockin scenario.

The connection handles, URLs, or credentials for these backing services or
attached resources should be obtained from the Configs. Any distinct services
that are attached to the application should be treated as an individual
resources.

![](https://cdn-images-1.medium.com/max/800/1*SWEanF1oRJLHwwzaewbKGw.png)
<span class="figcaption_hack">Application attaching the Database and Messaging service as backing service;
isolated form App</span>

In *Microservices and Kubernetes *paradigm, anything external to the service is
treated as an attached resource. There’s also an important behavior to note with
the containers and Kubernetes. When a container is deployed by default the
traffic is directed to it from the service. Kubernetes as an orchestrator
doesn’t know what are the backing services for that container by itself. But
when instructed, we can make Kubernetes wait for the container’s backing
services to be available before the traffic is routed in.

The developer's responsibility to check the list of backing services is
available and accessible for the container to fulfill the functionality when the
request is received. To achieve this, Kubernetes provides the *Probes* to make
use of.

> Probe describes a health check to be performed against a container to determine
> whether it is alive or ready to receive traffic.

There are three kinds of *probes *viz*., *`startupProbe`*, *`livenessProbe` and*
*`readinessProbe`* *explained below.

* `startupProbe` is used to determine whether the application/function within the
  container is started and ready for further process. The other probes are
  disabled until the *startupProbe is succeeded*. If this fails, the container
  will be restarted concerning the `restartPolicy` defined.
* `readinessProbe`** **is used to determine whether the container is used to
  determine if it is ready to service any requests.
* `livenessProbe`** **is used to determine whether the container is running or
  not. If fails, the `kubelet` will restart the container subject to
  `restartPolicy`.

The application should have a special method to validate the availability and
accessibility of all the backing services in an endpoint or a command. As these
can be used in the `readinessProbe` to determine the container is ready for
serving the requests.

More details on the probes can be found here on the [Kubernetes
site](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/).

### V. Build, release, run

> Strictly separate build and run stages.

Every application or service should have its own build, release pipelines. Since
all the services have their own repositories, it's very easy to have dedicated
pipelines for the repositories. There can be the use of common templates for the
pipelines. But not dedicated pipelines of the same code build for multiplatform.
Make use of the `Makefile` as much as possible to achieve the multiplatform
builds on the same code repository.

In our containers scenario, the build stage transforms the codebase into the
container image. During the build stage, all the vendor dependencies should be
taken care of and the system dependencies must be baked into the image. For
instance, while implementing the gRPC container’s health check, adding the
[grpc-health-probe](https://github.com/grpc-ecosystem/grpc-health-probe) binary
to the image during the build stage is recommended. Details of the gRPC health
check implementation can be found in this
[article](https://medium.com/@github.gkarthiks/implementing-healthchecks-in-grpc-containers-for-kubernetes-d5049989ab12).

![](https://cdn-images-1.medium.com/max/800/1*eF-WcHgqHbo1csnoG157og.png)
<span class="figcaption_hack">Build, Release and Run stage of a 12 Factor App</span>

Similarly, during the build stage, every build pipeline has a run id. The
container image must have either the code commit SHA or the build pipeline run
id embedded to it. And the code should be written in a way that this id is
logged along with the log messages, which helps in easily tracing back to the
source code that’s been built and running as a container in case of any issues.
For instance, in a GoLang container, this can be achieved by using the
`ldflags.` Running a `go build` command in the following way will overwrite the
version provided in the variable `Version.`

    go build -ldflags="-X main.Version=0.1.0" .

The release stage takes the build produced by the build stage and combines the
current config and makes it available for immediate execution in the execution
environment. For instance, the most predominant package manager that is used for
deploying the container images into the Kubernetes cluster is the
[helm](https://helm.sh). The release stage specifies that the helm should have
the corresponding files for the specific execution environment. In addition to
that, Helm’s post-rendering techniques allow the users with advanced
configuration needs to be able to use tools like `kustomize` to apply
configuration without a need to modify the chart.

The run stage (also known as “runtime”) runs the app in the execution
environment by launching some set of app processes against a selected release.
The run will have to be versioned. Like in the Helm chart, for example, there is
a release version associated with each run. This allows the application to be
rolled back to the previous version if needed. This enables the *Blue-Green
*deployment to have version N and N-1 available in production and easy to roll
back.

In Kubernetes, when an Operator is created, these run stages should be carefully
crafted to make use of the helm charts or its own custom way of maintaining the
release for rollbacks if necessary.

### VI. Processes

> Execute the app as one or more stateless processes.

The services should be stateless from the functional aspect, as there should be
a backing service read from the configuration and used for data persistence. In
the case of the containers, this is very much important as the containers are by
themselves ephemeral. They are prone to get restarted frequently so storing any
data within the container gets vanished. The 12-Factor app treats the processes
are stateless and shares nothing.

When dealing with microservices and containers, it's good to go with one
function per microservice per container. For instance, in a shopping website
example; the inventory service should be an individual service running in its
own container, and checkout service is another individual service running in its
own container and likewise, it goes on.

This allows to reduce the operational complexity and also reduces the *blast
radius* when an outage is happening. For instance, if the search service is
down, it doesn't affect the customers doing the checkout process as the checkout
service is running from a different container. Also enabling the services to
scale on demand. In this scenario, the checkout service is not going to face the
same load/traffic as the search service. So it is enough to scale up the search
service to serve the demand.

![](https://cdn-images-1.medium.com/max/800/1*F9O2NZE_Z7GrxbrOT9okzQ.png)
<span class="figcaption_hack">Search API has been scaled up independently to 2 instances</span>

Another most important best practice is running the service as a *non-root
process*. Often this is overlooked and that enables a very high-security risk.
Make sure to have the service run as a *non-root user* and provision that `uid`
in the image while building it.<br> <br> Also, while working with data from
within the containers, make sure the persistent data are never stored directly
inside a container, but on a corresponding storage volume or mount point
instead. And doing so, the newly introduced non-root uid mush have access to
those volumes to fulfill its functionality.

This can be achieved easily while building the container image. While creating
arbitrary user and user groups; to make the files and directories accessible for
this uid, these files and directories are chmod’ed. This will make the
Dockerfile slightly bloated and unnecessarily increasing the complexity of the
build process and size. The best practice is to use the existing least
privileged user in the base images. The two predominant base images `Ubuntu` and
`Alpine` have a common user called `nobody` with uid `65534`. So this should be
used for the creation of directories and as a run user.<br> <br> For example in
the `Dockerfile` do the following to use *nobody* effectively while creating
directories to permission *nobody*

    RUN mkdir -p /src/app
    RUN chown nobody:nobody /src/app

and to run the container as *nobody*

    USER nobody

In Kubernetes, if the container is running multiple processes, it means the
container needs resources viz., CPU, and memory to run those processes and the
below could happen. Ending up over utilizing resources from the cluster. If the
resource limits are specified for containers, this will end up abruptly
restarting the container after OOM Error by killing the running processes
whatsoever. So following the 12 Factor app and running one process per container
is really important.<br> <br> Also, you cannot guarantee that all the images
built are considering to run via non-root users. But in Kubernetes, using the
*Security Context*, a container can be restricted to execute to a particular
non-root user. This security context defines the privilege and access control
settings for a Pod/Container.<br> <br>

The security context is defined under the spec of a pod as shown in the below
snippet

If a `whoami` command is run on the container with the above-specified security
context, it will look like below

![](https://cdn-images-1.medium.com/max/800/1*g_rOw1auIQiy83KDO-G7hA.png)
<span class="figcaption_hack">examining user id with `whoami` command inside the container</span>

Now there are cases where more than one container runs in a Pod and might need
to run as different user IDs. In this scenario, a container-specific security
context can also be applied in the same way how the Pod Security context is
applied. An example of the container security context looks like

In this case, the container that has a specific context will be run as specified
and the rest of the containers will run as specified as in the pod level
context.

### VII. Port binding

> Use port binding to export services.

The traditional applications are deployed inside a webserver container. For
instance, once the application is built it will be deployed into Apache Tomcat,
GlassFish, WebLogic, etc., But with the modern and cloud-native applications,
they will be packaged with the webserver library as part of the application
itself.

In these scenarios, with microservices and containers, it's a best practice to
expose the service in the specific port number that is read from the environment
variable or configuration. This is classified as one of the *bootstrap
parameters* of the application. And another important aspect when designing a
containerized application is, that the service must listen for requests on
`0.0.0.0` and the port specified in the config and not on the `loopback`address.


Following this model improves the ease of portability for choosing between the
platform-as-a-service providers. Because of the Kubernetes’ built-in service
discovery model, these port-bindings can be abstracted with the mapping service
which is going to route the traffic to the containers.

In every language, there’s an elegant way of doing this. For example, in the
simple GoLang server, the below snippet shows how to read the port value from
the environment variable.

### VIII. Concurrency

> Scale-out via the process model.

The processes are the first-class citizens in the 12 Factor App. The
applications should be designed as multiple distributed processes that can
independently execute blocks of work and scale out by adding more processes.

The biggest advantage of having microservices and containers is scaling out the
individual service on demand. These individual services can be scaled without
disturbing any other services.

This type of scaling is called *Horizontal Scaling *as there are multiple
instances of the same service or application.* *Increasing the memory and CPU in
the same machine for the application’s usage is called *Vertical Scaling*. Refer
to the image in the VI. Processes for individual scaling of services. The
“Search API” which receives the heavy traffic has been scaled out to handle the
huge volume of incoming traffic.

This is also the reason the application or service should be stateless and no
data must be stored in the container.

In Kubernetes, the Horizontal Scaling is taken care of by an in-build monitor
called [Horizontal Pod
Autoscaler(HPA)](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/)*.
*The HPA automatically scales the number of pods based on the observed metrics.
The HPA can be configured for in-build CPU metrics or custom metrics as well
viz., number of requests per second, etc.,

> **Note**: One could use [Vertical Pod
> Autoscaler(VPA)](https://github.com/kubernetes/autoscaler/tree/master/vertical-pod-autoscaler#vertical-pod-autoscaler)
as well to increase the resources, but it doesn’t explain the concurrency factor
for the app.

A simple HPA manifest will look like the following

<span class="figcaption_hack">Sample Horizontal Pod Autoscaler</span>

Line numbers 6 and 7 states the maximum and minimum number of instances of the
pod selected with the references from lines numbers 8 to 11 when the given
metrics threshold which is CPU percentage of 75 on line number 12. Similarly,
the HPA can be also used to scale up and scale down with custom metrics as well.
For additional details, refer to the [Custome Metrics
API](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/instrumentation/custom-metrics-api.md).

When the Pods( or containers) are scaled up to multiple instances, each will
have its unique IP address allocated within the Kubernetes cluster. Despite they
have different IP addresses, the communication within the cluster is not
attached to the IP Addresses. These are abstracted by the Kubernetes services
with the match selector labels. These services will take care of routing the
traffic amongst these scaled pods.

> **Note**: HPA can’t be used for pods in *DaemonSets*. **💡**

The HPA feature of Kubernetes supports the *concurrency* of the 12Factor App.
But if all the scaled pods are deployed on the same worker nodes in the cluster
and if the worker node fails for some reason, the resiliency of the application
fails. To avoid this, Kubernetes offers another feature called *anti-affinity.
*There are multiple ways to schedule a pod to a node in Kubernetes. But the
anti-affinity feature provides a way not to schedule the pods on a node matching
certain conditions.

For example, additional instances of the pod can be restricted to be scheduled
on the same worker node if it already has an instance of the pod running in
that.

<span class="figcaption_hack">Pod Anti-Affinity spec</span>

In the above snippet of the YML manifest, the Pod Anti-Affinity is set to hard
rule and will not schedule another instance of the pod into the same Kubernetes
host. As the anti-affinity is set on the Kubernetes hostname at line number 5
with the matching labels on the pods on line numbers 8 and 9. This will make
sure to distribute the pods across the Kubernetes worker machines rather than
deploying all the pods on the same node.

### IX. Disposability

> The twelve-factor app’s processes are disposable, meaning they can be started or
> stopped at a moment’s notice.

The underlying infrastructure for the cloud-native application or service should
be always treated as disposable. This enforces the [Crash-only
design](http://lwn.net/Articles/191059/) for the application or services. This
means the application should handle the temporary loss of the underlying
infrastructure and/or sudden crash of the application or systems that form the
underlying infrastructure.

In such cases, the backing services should be isolated from the main
process/application. Refer to the *IV. Backing Services *for details on how to
isolate and what process to be isolated from the functional application. Also,
the configuration of the applications should be isolated from the main
application, refer to the *III. Configuration* factor for additional details on
it.

Another reason disposability is so critical is the scalability of an
application. A process that starts and shuts down fast allows for an environment
with auto-scalable application instances. The benefit of such a model is that
under higher load, more instances are created dynamically, and then as load
decreases, they are scaled-down. This allows serving the user needs without
major delays while also avoiding unnecessary costs by preventing overutilization
of resources.

Minimizing the startup time is more important to utilize the disposability
aspect. This is very much important in the case of the cron-jobs in the
Kubernetes environment and in any serverless custom-built runtime images.

The application should be designed to handle the graceful shutdown and restart.
This can be achieved by listening for the *SIGTERM* signal and doing the chore
works before the shutdown of the application. For example, if the application is
a web service, listening on the SIGTERM, and soon it's received the application
should stop receiving any new HTTP Requests and complete any pending previous
requests that are in process. In other cases, the application should be designed
to handle to release any resource lock when the SIGTERM signal is released
before shutting down the application. For example, any database connection
handlers, any record locks on the database, etc should be released before
shutdown.

When using microservices and containers, the principle of disposability is
implicitly followed to the maximum extent. There should be additional
considerations for isolating the request, session, and state from the
application containers. Using the robust message queues can enable this
furthermore. The usage of queueing technologies like Kafka, NATS, etc will
handle the unprocessed requests and tries to reconnect when the application
comes back up after an unexpected crash.

The below sample snippet in Golang describes how to handle the SIGTERM in a web
server code.

<span class="figcaption_hack">Handling SIGTERM in GoLang Web server</span>

In the above golang code snippet, at line number 12, a channel is made and is
listening on line 14. Then a new context is created with the Background function
for a specific time period. If there are any process needs to be completed, that
can be done in that time period and then the server is shutdown.

### X. Dev/prod parity

> Keep development, staging, and production as similar as possible.

The 12 Factor App methodology suggests keeping the gap between the Development
environment and Production environment as minimal as possible. Having
differences in the environments will lead to issues that are related to
incompatibility.

Designing and developing applications on the microservices architecture and
deploying via containers enables to minimize the differences amongst the
environment. Using the templated configurations makes it easier. Tools like
[Helm](https://helm.sh/), [Kustomize](https://kustomize.io/)*,
*[Kosko](https://kosko.dev/), etc helps in doing the templated configurations.

![](https://cdn-images-1.medium.com/max/800/1*qi8kmq3Bdx5jI2dyTVdiWQ.png)
<span class="figcaption_hack">Similar app deployed with templated configs via CI/CD</span>

*The twelve-factor developer resists the urge to use different backing services
between development and production.*** **Lightweight local services are less
compelling than they once were. Modern backing services such as Memcached,
PostgreSQL, and RabbitMQ are not difficult to install and run. Making use of
[Docker](https://www.docker.com/) and [Vagrant](https://www.vagrantup.com/) for
these technologies helps in running those services locally which closely
approximate production environments. The cost of installing and using these
systems is low compared to the benefit of dev/prod parity and continuous
deployment.

For the Kubernetes clusters, there are numerous local clusters tools are
available. Tools like [Kind](https://kind.sigs.k8s.io/),
[minikube](https://minikube.sigs.k8s.io/docs/),
[microK8s](https://microk8s.io/), etc are available to simulate the production
Kubernetes cluster. These clusters should be utilized to deploy the backing
services as well to bring the maximum resemblance to the Production environment.
Do so in the development environment will minimize the issues caused by the
compatibility between the technologies.

### XI. Logs

> Treat logs as event streams.

Logs provide paramount information to the behavior of a running application.
Especially in helping troubleshoot the issues that may arise while running the
applications or service. The logs are treated as critical data streams triggered
when an event occurs. The applications should be designed to log the information
in `stdout`, it should not attempt to write to or manage logfiles. A
twelve-factor app never concerns itself with routing or storage of its output
stream, it is unaware of the logging infrastructure.

Logging should be always decoupled from the application and is very useful when
the apps require dynamic scaling. Because of the dynamic scaling nature, when
the logging infrastructure is decoupled, it reduces the overhead of managing the
log storage location.

There are a lot of tools that can further collate the logs from the application
and can be used for analysis. Very rich log information can be further used to
study and alert the process owners proactively about a crash that’s going to
happen.

In containers, it's very important to have a structured log. As the containers
are built and deployed via CI/CD in a matter of hours, if possible the version
of the application built should also be logged as part of the log information.
This helps in reducing the time taken to backtrace the code commit version when
an issue is reported in a running application. Refer to the “*V. Build, release,
run”* factor for `ldflags` a toolchain used in GoLang language for injecting the
build version into the container images.

### XII. Admin processes

> Run admin/management tasks as one-off processes.

The administrative tasks are one-off tasks as per 12 Factor App. For the cloud
native applications, these can be a one-off task or timed repetable tasks like
running reports, taking backups of the databases, etc,. This section will
discuss about those timed jobs for the cloud-native applications.

Often timed jobs are executed as cron and handled by the application by
themselves. This introduces the tight coupling of the cron logic with the
application logic. Also increases the maintainence overhead for the dynamically
scaling apps and very hig cost maintenance if these apps are across timezones.

It is always recommended to do these kind of timed logic as a separate
application by itself. That follows all the 12 Factor methodologies as well.

In Kubernetes, these are executed as *CronJobs*. CronJobs are first class
citizens of the Kubernetes resources. This follows the Cron schedule to get
executed. A sample CronJob YAML manifest will look like the below.

The above cronjob executes at every hour and prints the date. The *schedule* at
line number 6 is used for the run schedule and follows the cron syntax. It is
important to make the container image of the CronJob really light weight. There
are additional known limitations to cron job and can be found
[here](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/#cron-job-limitations).

The one-off admin tasks such as warming the cache, bootstrapping the database,
etc., can  also be handled within the application. But as stated above these
increases the tight coupling of the admin task’s logic with the app logic. These
can be isolated and deployed as a separate specialized containers that run
before app containers in a pod. These are called as *Init containers.*

Init Containers can contain utilities or setup scripts not present in an app
image. Following the similar 12 Factor App methodology, the init containers are
built as as application by itself. During the release stage of the application,
these init-containers are added to the run. An app can have multiple one-off
startup process. Similarly there can be any number of init-containers. The order
of the init-containers in the manifest determines the order of execution.

The above snippet of code shows the init container that’s waiting for the
database to be available. The main application container will not start untill
the init containers ran successfully as that’s the order of execution.

*****

Following 12 Factor App principles for developing cloud-native application will
result in more standard, robust, efficient, portable applications that uses
declarative formats, suitable for deployment on any platforms, minimize
difference between development and production environment and can scale on
demand.

But the twelve-factor methodologies may not be the silver bullet for solving
distributed computing. Also, if you are trying to overcome the burden of legacy,
on-premise applications, 12 factor is not ready for you.

*Whats Next?* Check out the [Application
Continuum](https://www.appcontinuum.io/) to understand more on cloud-native
architectire. 
