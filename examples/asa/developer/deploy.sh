currentNS=$(kubectl config view --minify --output 'jsonpath={..namespace}')
rm -r $currentNS
mkdir $currentNS
cd $currentNS
sed s/currentNS/$currentNS/g ../*.yaml > temp.yaml
kubectl apply -f temp.yaml
cd ..
