currentNS=$(kubectl config view --minify --output 'jsonpath={..namespace}')
rm -r ../temp/$currentNS
mkdir ../temp
mkdir ../temp/$currentNS
cd ../temp/$currentNS
sed s/currentNS/$currentNS/g ../../developer/*.yaml > temp.yaml
kubectl apply -f temp.yaml
cd ../../developer
