kubectl delete -f git-source-template.yaml

kubectl delete -f kpack-image-template.yaml
kubectl delete -f kpack-image-template-directjar.yaml

kubectl delete -f app-deploy-containerapp-cluster-template-noconfig.yaml
kubectl delete -f app-deploy-template.yaml


kubectl delete -f supply-chain-containerapp-directjar.yaml
kubectl delete -f supply-chain-jar.yaml
kubectl delete -f supply-chain.yaml