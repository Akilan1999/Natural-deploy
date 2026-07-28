# Natural Deploy

## Motivation:
Have you ever felt when using any declarative type of program that is used for deployment. Its really unnatural to use (i.e you are a YAML developer). This is in comparison to using any programming language (Which has advanced static type checkers, error diagnosis etc…).

The 2nd issue is it’s really hard to setup servers and can only work with nodes which have public IPs. Most orchestration tools are built on the basis to use static IPs and fixed ports.

The effect of this move is that you can easily get vendor locked and you are less likely to choose the economical choice of owning your own hardware. We have lesser gratefullness of how complex problems were abstracted in simple tools like FTP, GIT, OrgMode or the design principles of Unix utilities.

The aim of “Natural deploy” is to be library build on top of P2PRC to make it as easy as running a binary to deloy your custom tasks to your servers. This will gives you the flexibility to truly customise building your orchestration layer to run your programs across your cluster of servers which can be situated anywhere.

## The Solution:
<img width="778" height="834" alt="Proposal" src="https://github.com/user-attachments/assets/4ce16823-3be2-4b89-924c-94331e9a627e" />

## Documentation
- https://github.com/Akilan1999/Natural-deploy/blob/master/docs/index.org
