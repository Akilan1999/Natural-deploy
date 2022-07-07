# Natural Deploy 

Work in progress: Will continue from August 2022. 

Its Go way of doing Ansibles:

## Motivation:
Have you ever felt when using ansible or any declarative type of program that is used for deployment. Its really unnatural to use. This is in comparison to using any programming language. And running ansibles is another nightmare. 
More on Anibles cannot deploy on windows natively. There are more problems and my list goes on ... 

## The Solution: 
Programming languages are powerful to act as an alternative to ansibles. 
I plan is to use a golang to run tasks on other computer and embed all the other files to the binary. So when running the binary it can unpack all the files and could run on the computer. Since this approach is pure Go its run on all architectures and platforms supported by golang. 

The implementation to run them across different platforms is also pure Go and can could either use SSH to connect to the other machines or whatever windows uses for remote connection. 
