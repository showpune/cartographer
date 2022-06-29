# secret to write the customer image
kubectl create secret docker-registry private-registry-credentials --docker-server=https://index.docker.io/v1/ --docker-username=showpune --docker-password=Nuonuo8815 --docker-email=zhiyongli@microsoft.com 

# secret to read the shared builder
kubectl create secret docker-registry shared-registry-credentials --docker-server=serviceacr.azurecr.io --docker-username=serviceacr --docker-password=1dXPkDXTRv+DlZxi/n8HRC0IjVdKRadY --docker-email=zhiyongli@microsoft.com