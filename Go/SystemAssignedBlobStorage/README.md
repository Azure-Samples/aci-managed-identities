# Access Blob Storage from Azure Container Instances Using Manged Identites

This example walksthrough how to access an Azure Blob Storage from Azure Container Instances. The best part of all, absolutly no access credentials needed.

## Prerequisites

- An Azure subscription. Get a [free trial here](https://azure.microsoft.com/en-us/free/)!

- A [DockerHub](http://dockerhub.com) or other container registry account.

- A machine with Docker installed

## Getting Started

### Build and Push the Image

First up is to build and push our container to a container registry e.g. Dockerhub, Azure Container Registry

1. On a machine with docker installed, run the following command in this directory with the Dockerfile to build the container image

```sh
    docker build -t <dockerhub-username>/msi-blob:0.0.1 .
```

The above command will install all the dependencies into the container. The first time you run this is will take a while to download all of the dependencies but will be much faster once it is cached.

2. Push the image to [DockerHub](http://dockerhub.com)

```sh
    docker push <dockerhub-username>/msi-msi-blob:0.0.1
```

### Setting Up the Azure Resources

To deploy the container, open either a terminal with the Azure CLI installed or check out [CloudShell](https://shell.azure.com/) for a ready-to-go option

First, lets set up some environment variables to make these commands nicer to copy and paste

```sh
DOCKER_IMAGE_NAME="<imageName>" #This is the image name you pushed to dockerhub
RESOURCE_GROUP="<myResourceGroup>" #If this doesn't exist we will create one
STORAGE_ACCOUNT="<storageaccountname>" #This must be all lowercase or numbers, no special characters
```

Create a resource group.

    az group create --name $RESOURCE_GROUP --location westus

### Setup the Azure Storage Account

1. Create the storage account

    az storage account create -g $RESOURCE_GROUP -n $STORAGE_ACCOUNT

2. Let's save the storage account ID and Key for later

    STORAGE_KEY=$(az storage account keys list -g $RESOURCE_GROUP -n $STORAGE_ACCOUNT --query [0].value | tr -d '"')

    STORAGE_ACCOUNT_ID=$(az storage account show -g $RESOURCE_GROUP -n $STORAGE_ACCOUNT --query id | tr -d '"')

3. Create a blob container inside the storage account

    az storage container create -n testdata --account-name $STORAGE_ACCOUNT --account-key $STORAGE_KEY

4. Upload a simple file to the blob. If you don't have a text file handy simply run: `echo "hello from blob storage!" > ./testfile.txt`

    az storage blob upload -c testdata -n testfile.txt -f ./testfile.txt --account-name $STORAGE_ACCOUNT --account-key $STORAGE_KEY

### Time to Deploy to Azure Container Instances

Now we can use the single deploy command for ACI, make sure to change to the correct image name below:

```sh
az container create \
    --resource-group $RESOURCE_GROUP \
    --name msi-blob \
    -e STORAGE_ACCOUNT_ID=$STORAGE_ACCOUNT_ID \
    --image $DOCKER_IMAGE_NAME \
    --assign-identity --scope $STORAGE_ACCOUNT_ID
```

The above command will create the container instance and set up everything needed for the system assigned Managed Identities.

Finally, once the command finished, check the log output to see the text file:

    az container logs -g $RESOURCE_GROUP -n msi-blob

That should be a good start to never needing to store production credentials again.

## Issues

If you have any issues or find any mistakes, Please open an Issue on this repository and we will update this document.
