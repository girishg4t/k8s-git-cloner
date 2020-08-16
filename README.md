This project is to demonstrate the cloning of git repo inside a container, it uses skaffold templete to create docker file and deploy the application into kubernetes , this application does the following steps

1) Requests Github for Read-access on the repositories 
2) Github asks for user credentials, user enters credentials
3) Github authorizes and gives token to App, 
4) The App clones the repository on behalf of the user. 

You can run the application using 

```skaffold dev```

for more details on skaffold visit https://skaffold.dev/docs/quickstart/