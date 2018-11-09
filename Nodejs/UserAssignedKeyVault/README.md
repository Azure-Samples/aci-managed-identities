# Access KeyVault from Azure Container Instances Using Manged Identites

This example walks through how to get a secret from key vault. But no more passing tokens or credentials to your application to access key vault. We're going to be using Managed Identites for Azure Resources instead.

Our sample will retrieve the secret from key vault and log it.

## Prerequisites

- An Azure subscription. Get a [free trial here](https://azure.microsoft.com/en-us/free/)!

- A [DockerHub](http://dockerhub.com) or other container registry account.

- A machine with Docker installed

## Getting Started

### Build and Push the Image

First up is to build and push our container to a container registry e.g. Dockerhub, Azure Container Registry

1. On a machine with docker installed, run the following command in this directory with the Dockerfile to build the container image

```sh
    docker build -t <dockerhub-username>/msi-nodejs:0.0.1 .
```

The above command will install all the dependencies into the container. The first time you run this is will take a while to download all of the dependencies but will be much faster once it is cached.

2. Push the image to [DockerHub](http://dockerhub.com)

```sh
    docker push <dockerhub-username>/msi-nodejs:0.0.1
```

### Setting Up the Azure Resources

To deploy the container, open either a terminal with the Azure CLI installed or check out [CloudShell](https://shell.azure.com/) for a ready-to-go option

First, lets set up some environment variables to make these commands nicer to copy and paste

```sh
DOCKER_IMAGE_NAME="<imageName>" #This is the image name you pushed to dockerhub
RESOURCE_GROUP="<myResourceGroup>" #If this doesn't exist we will create one
KEYVAULT_NAME="<mykeyvault>"
SECRET_NAME="<secretName>"
SECRET_VERSION="secretVersion>"
USER_ASSIGNED_IDENTITY_NAME="<myUserAssignedIdenity>"
```

### Create a resource group.

    az group create --name $RESOURCE_GROUP --location westus

#### Create the User Assigned Identity

Next we create a user assigned identity. Once we create and add the permissions, we will be able to use this for multiple container groups.

The following commands will create the identity and get the needed information from it.

    CLIENT_ID=$(az identity create -g $RESOURCE_GROUP --name $USER_ASSIGNED_IDENTITY_NAME --query clientId | tr -d '"')
    PRINCIPAL_ID=$(az identity show -g $RESOURCE_GROUP --name $USER_ASSIGNED_IDENTITY_NAME --query principalId | tr -d '"')
    MSI_RESOURCE_ID=$(az identity show -g $RESOURCE_GROUP --name $USER_ASSIGNED_IDENTITY_NAME --query id | tr -d '"')

#### Create a Key Vault and Set the Permissions

If you don’t already have a Key Vault create, use the following command to create one:

    az keyvault create -g $RESOURCE_GROUP --name $KEYVAULT_NAME

Now we create a secret in KeyVault

    SECRET_VERSION=$(az keyvault secret set --name $SECRET_NAME --value <secret value> --vault-name $KEYVAULT_NAME | cut -d'/' -f6 | tr -d '"')

Now, we can give our identity access to the Key Vault

    az keyvault set-policy -n $KEYVAULT_NAME --object-id $PRINCIPAL_ID -g $RESOURCE_GROUP --secret-permissions get

The above command uses the environment variable we set to give our identity “get” permission for secrets in the Key Vault

### Time to Deploy to Azure Container Instances

Now we can use the single deploy command for ACI, make sure to change to the correct image name below:

```sh
az container create \
    --resource-group $RESOURCE_GROUP \
    --name msi-nodejs \
    -e KEY_VAULT_NAME=$KEY_VAULT_NAME SECRET_NAME=$SECRET_NAME SECRET_VERSION=$SECRET_VERSION \
    --image $DOCKER_IMAGE_NAME \
    --assign-identity $MSI_RESOURCE_ID
```

The above command will create the container instance and set up everything needed for Managed Identities. 

Once the command has finished, we can run the below command to get the container logs which output the secret value : 

```sh
    az container logs --resource-group $RESOURCE_GROUP --name msi-nodejs
```
The result will be: "The secret value is: XXXX"

## Issues

If you have any issues or find any mistakes, Please open an Issue on this repository and we will update this document.
