kubectl apply -f git-source-template.yaml

kubectl apply -f kpack-image-template.yaml
kubectl apply -f kpack-image-template-directjar.yaml

kubectl apply -f app-deploy-containerapp-cluster-template-noconfig.yaml
kubectl apply -f app-deploy-template.yaml


kubectl apply -f supply-chain-containerapp-directjar.yaml
kubectl apply -f supply-chain-jar.yaml
kubectl apply -f supply-chain.yaml