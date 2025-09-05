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
