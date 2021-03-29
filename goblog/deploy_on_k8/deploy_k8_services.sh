#!/bin/bash
kubectl apply -f acc_deployment.yaml # create the accountservice pods as a k8 deployment
kubectl apply -f acc_service.yaml # use the accountservice deployment to create the service

kubectl apply -f quotes_deployment.yaml # quotes deployment
kubectl apply -f quotes_service.yaml # service based on the quotes deployment