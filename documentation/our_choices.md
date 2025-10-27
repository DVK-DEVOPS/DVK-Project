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

After you have setup a these code quality tools and gone through the issues, your group should create a brief document that answers the following questions:

# Do you agree with the findings?

# Which ones did you fix?

# Which ones did you ignore?

# Why?
