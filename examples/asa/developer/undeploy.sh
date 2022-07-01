currentNS=$(kubectl config view --minify --output 'jsonpath={..namespace}')
read -e -p "Enter Your namespace:" -i $currentNS currentNS
kubectl delete workload -l app.tanzu.vmware.com/workload-type=web --namespace=$currentNS
kubectl delete workload -l app.tanzu.vmware.com/workload-type=web-directjar --namespace=$currentNS
# kubectl delete namespace $currentNS
