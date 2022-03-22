# Compile the echoserver project to create a Linux binary for use in Kubernetes.
go_compile = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/linux/echoserver'
local_resource(
    'echoserver-compile',
    go_compile,
    deps=[
        'go.mod',
        'go.sum',
        './main.go',
        './cmd',
        './internal',
        './pkg'
    ]
)

# Build the Docker image of echoserver to send to the k3d registry for Kubernetes to pull from.
# This is synchronized to the compile of the Linux binary.
docker_build(
    'echoserver', '.',
    dockerfile='infrastructure/local/deployment/Dockerfile',
    only=[
        './build/linux'
    ],
    live_update=[
        sync('./build/linux', '/build')
    ]
)

# Use the Deployment and Service manifests to create the definitions for the echoserver deployment.
k8s_yaml('infrastructure/local/deployment/k8s.yaml')

# Push echoserver as a Kubernetes resource when the complitation step completes.
k8s_resource(
    'echoserver',
    resource_deps=[
        'echoserver-compile'
    ]
)
