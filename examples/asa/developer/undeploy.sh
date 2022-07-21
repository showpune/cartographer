currentNS=$(kubectl config view --minify --output 'jsonpath={..namespace}')
read -e -p "Enter Your namespace:" -i $currentNS currentNS
kubectl config set-context --current --namespace=$currentNS
kubectl delete workload -l app.tanzu.vmware.com/workload-type=web 
kubectl delete workload -l app.tanzu.vmware.com/workload-type=web-directjar
kubectl delete workload -l app.tanzu.vmware.com/workload-type=containerapp
kubectl delete workload -l app.tanzu.vmware.com/workload-type=containerapp-directjar

items=$(kubectl get pods,App --output 'jsonpath={.items}')
while [ "$items" != "[]" ] ;do
echo 'There is alive pods, wait for them to be deleted'
sleep 5
items=$(kubectl get pods,App --output 'jsonpath={.items}')
done

echo 'Pod/App is cleared, now delete the namespace'
read -e -p "Delete namespace:" -i 'N' isDelete
if [ "$isDelete" = "Y" ] ;then
kubectl delete namespace $currentNS
fi