load('ext://restart_process', 'docker_build_with_restart')

go_compile = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/linux/echoserver'
local_resource(
    'echoserver-go-compile',
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

docker_build_with_restart(
    'echoserver', '.',
    dockerfile='infrastructure/local/deployment/Dockerfile',
    entrypoint='/build/linux/echoserver',
    only=[
        './build/linux'
    ],
    live_update=[
        sync('./build/linux', '/build')
    ]
)

k8s_yaml('infrastructure/local/deployment/k8s.yaml')
k8s_resource(
    'echoserver',
    resource_deps=[
        'echoserver-go-compile'
    ]
)
