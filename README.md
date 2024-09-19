# plateforme-mycli

## Initialisation du projet

Pour installer la CLI, exécutez simplement :

```
$ go get github.com/urfave/cli
```

Assurez-vous que vous `PATH`incluez le `$GOPATH/bin`répertoire afin que vos commandes puissent être facilement utilisées :

1jouter le répertoire `$GOPATH/bin` à votre variable d'environnement `PATH`. Cela permet à votre système de trouver les exécutables Go que vous avez installés. Voici comment vous pouvez le faire sur une machine Windows :

1. Ouvrez l'invite de commandes ou PowerShell.
2. Exécutez la commande suivante pour ajouter `$GOPATH/bin` à votre `PATH` pour la session en cours :

```
export PATH=$PATH:$GOPATH/bin
```

### Plateformes prises en charge

[](https://github.com/minio/cli#supported-platforms)

cli est testé sur plusieurs versions de Go sous Linux et sur la dernière version publiée de Go sous OS X et Windows. Pour plus de détails, voir [`./.travis.yml`](https://github.com/minio/cli/blob/master/.travis.yml)et [`./appveyor.yml`](https://github.com/minio/cli/blob/master/appveyor.yml).
