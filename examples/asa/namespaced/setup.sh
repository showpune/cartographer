read -p "please enter your namespace: " namespace
kubectl create namespace $namespace
kubectl config set-context --current --namespace=$namespace
export kubenamespace=$namespace

read -s -p "please enter your private password: " dockerpassword
echo 'Create secret to write the customer image'
kubectl create secret docker-registry private-registry-credentials --docker-server=https://index.docker.io/v1/ --docker-username=showpune --docker-password=$dockerpassword --docker-email=zhiyongli@microsoft.com 

echo 'Create secret to read the cluster resource'
kubectl create secret docker-registry shared-registry-credentials --docker-server=serviceacr.azurecr.io --docker-username=serviceacr --docker-password=1dXPkDXTRv+DlZxi/n8HRC0IjVdKRadY --docker-email=zhiyongli@microsoft.com

read -p "Input the subscription id for container app: " subscriptionid
read -p "Input the resourcegroup for container app: " resourcegroup
read -p "Input the ask name for container app: " aksname
az account set --subscription $subscriptionid
az aks get-credentials --resource-group $resourcegroup --name $aksname --admin --overwrite-existing --file ./containerapp.yaml

echo 'Create secret to access customer aks'
kubectl create secret generic customer-aks --from-file=kubeconfig=./containerapp.yaml

echo 'Create namespace sa, role and rolebinding'
kubectl apply -f service-account.yaml

