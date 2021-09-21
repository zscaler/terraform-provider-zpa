
export_cred(){
    export ZPA_CLIENT_ID="MjE2MTk2MjU3MzMxMjgyMDcwLTg0MTgyN2I5LTEwMDQtNDE0Mi1iYjQwLTVlOGE0NWEyMjc2MQ=="
    export ZPA_CLIENT_SECRET="HBRM'}IQgum#Yd~VxDz*d]@X6]Zab)<N"
    export ZPA_CUSTOMER_ID="216196257331281920"
    export TF_LOG=TRACE
}

ta(){
    export_cred
    cd "./$1"
    rm -rf .terraform*
    terraform init && terraform apply --auto-approve > init.log 2>&1
    cd "../"
}

td(){
    export_cred
    cd "./$1"
    terraform destroy --auto-approve  > delete.log 2>&1
    cd "../"
}




