# Choices for the rewrite
Below you will find the choices we have picked for rewriting the legacy codebase and the reasoning behind. 

## Why? 
We realise that DevOps is about promoting a culture of collaboration and shared responsibility that will ultimately result in a better product both in quality and in delivery. By exposing our team culture in a markdown file like this serves to publicly declare our intents and reasoning both for readers visiting our repository but also for ourselves; think of it as a public document that the team can go back and visit in order to "remember" why we made the choices we did. 


As stated above putting our thoughts into a public document in a manner such as this also ensures that everyone on the team is on the same page. This can help alleviate the situations where a team member might ask themselves "Why did we choose Go" and "why didn't we just split our codebase into two repositories consisting of a back- and frontend?". 

This document tries to establish transparency about the decision making of our team helping to ensure collective ownership of the system.

# Multi- og monorepo? 
Despite having undertaken several projects where we've split up our code in a front- and backend repository we've gone with a monorepo structure for this rewrite project.


`It is our current understanding that by going with the monorepo structure we simplify our system and our processes.`

The team believes that a project of this size is perfect for a legacy system of this size. The legacy project contains around 220 lines of code which obviously is a small system. At the moment of writing we do not have the expectation of writing a system drastically larger than the original.

The team believes that it will be easier to configure a monorepo with the future Docker configuration and in configuring the CI/CD pipeline.

The team believes that it will be easier to debug potential issues within the system if we keep the codebase within a single repository. This is further helped by preventing our developers from having to swap between two repositories or projects which will help with the efficiency of tracing potential issues.


# Go & Gorilla
`Go & Gorilla presents a fun challenge and learning opportunity.`

The team has exclusively worked with Java and Spring in previous projects which meant that the team was motivated to pick Go & Gorilla for its' minimalistic and highly readable syntax. 

Go is highly suited for small scale systems such as the legacy system that the team is trying to rewrite.


# SQLite Database
`Simple & Portable`

As has been repeatedly stated above the scale of the legacy project inevitably plays a role in going with an SQLite database. The team wanted *simplicity*, *portability* and *low overhead* which an SQLite database provides. SQLite provides fast performance and is a lightweight solution. Perfect for our project.

# How we work with environment variables

Azure Vault + env variables in systemd config on the server.

# What version control strategy we chose and how we enforced it?

We decided to use GitHub Flow. For every new feature or bug fix, we create a separate feature branch with a descriptive name. We regularly commit and push our work to the remote.

Before merging to main, we make a pull request, so the code is reviewed by at least one team member (we have set up branch protection rules on GitHub so no one can push directly to main). The point of the review is rather informational than quality control - it is that everybody knows exactly what the other team members have been working on. Therefore, we do not write an actual review on GitHub but simply merge the pull request when it was created. In that regard our flow resembles Feature branching strategy.
If someone has already merged new changes into main, we pull those updates into our feature branch before creating the PR to avoid conflicts. Once a pull request is approved and merged, the code can be deployed right away.

So overall, we only have one long-living branch (main), and the feature branches are merged frequently.

We chose git merge as our branch integration strategy because we want to preserve the commit history the way it actually happened and clearly see when each feature branch was merged into main.

# Why did we choose the one we did? Why didn't we choose others?

The reason why we chose GitHub flow is that it’s a simple branching strategy that fits well with our team size and the project which is not very complex.  The alternative to GitHub flow was Git flow and Trunk based development.

Git Flow was an overkill for our small team, as it uses long-living branches (develop and release), which would make things more complicated and slow down our delivery.

Trunk based development seemed too risky, since it means committing directly to main. Given that we have no tests yet, it would increase the chance of bugs in production.

# What advantages and disadvantages did we run into during the course?

One of the advantages we experienced is that the workflow is familiar and simple, so there haven’t been any merge conflicts.  Another thing we enjoyed about this strategy is that it works nicely with our CI/CD pipeline, and the main branch always stays stable and ready for deployment.

The disadvantage we encountered was that even small fixes must go through a pull request, which slows down the process. Also, the pull request templates can feel a bit annoying to fill in for quick changes.

# How are you DevOps?

At the start of the project, we wrote down our motivation for the rewrite publicly on our repository. Furthermore, we noted down the conventions we aimed to adhere to during the rewrite. Starting out, we noted some of the problems with the codebase that we needed to rectify in our rewrite project.
We have also tried to utilise the wiki section of the Github repository in regard to documenting a potential database migration tool with an accompanying branch and ORM tools compatible with Golang.

*We are DevOps...*
Our group has on several occasions utilised swarming in order to deal with a critical issue: we had an issue with our push-based deployment creating weird folder structures and not running the binary on our VM server. Physically we sat and troubleshooted the issue together until it no longer created the odd folder structures and ran the correct binary.

*Culture of collaboration and shared responsibility, collaboration and communication & transparency, visibility, and knowledge sharing!*
There’s been a culture of collaboration in our group in the sense that we needed to work together in order to overcome the workload that this rewrite project possesses while being able to maintain our other electives and other commitments on the side. Each member of the group wants the others to succeed, which ties into communication within the group. The group has worked together before and knows each other, which has no doubt influenced the way that we communicate with each other. 
Presumably because the group has worked with each other before communication has also been permeated with kindness and the willingness to ask questions like “How does this work?” without experiencing judgement.

The group ensures psychological safety between members to ask questions which the group also recognizes as a different way of troubleshooting potential issues. Psychological safety encourages members to voice ideas and uncertainties freely. We believe that this greatly improves the group’s combined ability to troubleshoot problems from different perspectives. The goal is to minimize self-doubt but also ensure that valuable ideas and input are not held back to the detriment of the entire project.

The group has strived for descriptive commit messages in regard to larger changes but as complexity has grown within certain implementations commit messages have degraded from sentences to single words when troubleshooting problems with small or incremental changes to the respective files. However, with the use of Github Copilot the changes a pull request seeks to implement is transformed into informative and descriptive messages that allow the other members of the group to see the impact that a pull request will make to the branch that the pull request targets.

Many problems and “issues” were added to the Github Projects board in order to give the group an overview of pending tasks to fulfill the weekly DevOps elective requirements for the rewrite project. The DevOps elective introduced several new features and external improvements that we had to prioritise alongside the ordinary rewrite. 

*The A in CALMS & Fail fast, recover fast*

Our group has managed to automate a great deal over the current lifecycle of the project. As of writing, the project repository currently has 11 automated Github Action workflows. We have automated linting which inspects the project code and the accompanying Dockerfiles which helps each member follow consistent styling resulting in cleaner and more uniform code within the project.
We have automated both Continuous Delivery and Deployment, which ensures that whenever new changes are introduced to the main branch of  the project repository, these changes are reflected in the live environment within minutes. 
The group has predominantly utilised pull requests in order to merge feature branches with the main branch. This resulted in a “long” and frustrating period of time before a member of the group could get their changes merged because we chose to restrict direct pushes to the main branch. This meant that every member that wanted to push changes had to wait for approval if anyone was even available. Recently, the group has automated the pull request approval so that the Github Actions bot will approve any pull request made by named members of the group. However, the group recognises that this approach necessitates a large degree of mutual trust and self-regulation in order to work properly and that it may not fit the “correct” approach to collaboration, but due to lots of other activities the group decided to make this workaround.
The group has therefore imposed upon themselves the convention that any substantial changes to the overall project have to be approved by another member of the group and not solely the auto approval workflow.
 
*We are NOT DevOps...*
During the rewrite project complexity has grown tremendously with the addition of logging and monitoring. As we are near the end of the semester and the DevOps elective and the accompanying workload, it is inevitable that the group has been faced with burnout. This means that knowledge sharing between members of the group has been limited since the developer usually just wants to finish the feature and implementation rather than taking more time out of their day to share their work and writing up documentation.  In practice, we created avoidable problems for ourselves by pushing forward without the time or energy to coordinate effectively, and by that limiting our ability to be more effective devops. 

# Software Quality
*Do you agree with the findings?*
In one instance SonarQube aided in ensuring security in a sensitive user flow where the user can reset their password. SonarQube flagged the code creating a cookie for the user in order to store information.
*Which ones did you fix?*
As SonarQube flagged the cookie creation code we believed it prudent to make sure that the cookie should not be sent over plain HTTP connections as this could result in a man-in-the-middle-attack where an attacker could read the cookie value and possibly hijack the password reset flow. SonarQube suggested we adjust our code, so the cookie was no longer sent over HTTP but exclusively on HTTPS, so the cookie was encrypted in traffic.
*Which ones did you ignore?*
Some lines of error handling within the Go code were flagged by SonarQube but were ultimately ignored by the group.
*Why?*
These lines of error handling were rewritten to the current structure in order to be compliant with the golangci linter. The SonarQube flag of said code was not deemed to be critical enough to change as SonarQube describes it as being an “unnecessary variable declaration”.
 # <img width="468" height="107" alt="image" src="https://github.com/user-attachments/assets/c1d00712-6c32-4911-a75c-ad98923a5691" />


# Monitoring realization
 <img width="504" height="278" alt="image" src="https://github.com/user-attachments/assets/a6f1d5a0-05fa-49dc-bcc2-be33a43eb51b" />

Among other things, we monitored RAM usage per container on both the production VM and the monitoring VM. We chose to monitor this because we knew our resources on Azure were limited in terms of RAM and CPU. The monitoring helped us to identify which containers consumed the most resources. It was useful later when we needed to decide which technologies to remove from the monitoring setup, in order to reduce latency on the VMs without purchasing additional capacity on Azure. On the dashboard we saw that cAdvisor’s RAM load was one of the heaviest, so we removed it from the monitoring VM.
Another aspect we monitored was the search words users entered when using the system’s search feature. Understanding user preferences helps tailoring the system to provide more relevant data, making successful searches more frequent and therefore attracting more users. We can now populate the database by scraping pages containing information relevant to the searches we observed most often (fx Netflix page for searches related to movies).
 <img width="526" height="271" alt="image" src="https://github.com/user-attachments/assets/1f245a36-6e5d-485f-8675-c231d2f653bc" />

And we would like to mention another insight which came as a result of the whole logging-monitoring setup. Our production VM was constantly trying to restart a failing systemd service which was used to run our app before we transitioned to dockerization (we saw it in the syslog dashboard). Before the central logging had been established we were not aware the service hadn’t been properly stopped. It has now been removed.

