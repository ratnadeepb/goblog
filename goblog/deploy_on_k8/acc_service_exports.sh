#!/bin/bash
export SVC_HOST=$(kubectl get service/accountservice -o jsonpath='{.spec.clusterIP}')

export SVC_PORT=6767