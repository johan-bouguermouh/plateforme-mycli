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

[Voir ici](https://github.com/minio/cli#supported-platforms)

cli est testé sur plusieurs versions de Go sous Linux et sur la dernière version publiée de Go sous OS X et Windows. Pour plus de détails, voir [`./.travis.yml`](https://github.com/minio/cli/blob/master/.travis.yml)et [`./appveyor.yml`](https://github.com/minio/cli/blob/master/appveyor.yml).

### Utilisation d'une Image minio pour les tests

Lancer Docker Desktop

```bash
// A la racine du projet
docker-compose up --build
```

Cela lancera automatiquement l'image docker en local sur le port `9001`. Pensez à changer le nom `MINIO_ROOT_USER` & `MINIO_ROOT_PASSWORD` par vos propres credentials.

## Usage du CLI

### Avant Propos

Cette version du CLI comprend une authentification automatique en `HMAC-SHA256`. Par conséquent il n'est pour l'instant que compatible avec la version4 d'authentification des buckets Amazon. [Voir la documentation officielle](https://docs.aws.amazon.com/fr_fr/AmazonSimpleDB/latest/DeveloperGuide/HMACAuth.html)

Nous prendrons comme exemple cette image docker pour rentré dans le détaille d'utilisation du CLI

### Usage des Alias

La command alias comprend à ce stade trois commandes :

- set
- list
- current

#### Alias Set

La command Alias set permet de créer un alias, cela facilitera par la suite la manière dont vous essayez de rentrez en contact avec votre bucket S3.

Cette commande prend comme première argument le nom de votre alias

Il comprend les Flags suivant :

- **`-port, -p`** **(obligatoire)**: Définit le port de votre server S3- \*-hostname,
- _`-hostname, -H`_ **(obligatoire)**: Définit l'hôte auquel nous devons nous connecter
- _`-keyname, -k_ `**(obligatoire)**: Définit le AWSAccessKeyId
- `_-Secret, -s`\_ **(obligatoire)**: Votre clef secret qui sera ensuite signée en HMAC-SHA
- `_-current, -c_` **(Optionnel)**: Permet de signifier directement au programme que cette alias doit être pris par défault

##### Exemple d'utilisation de Alias Set

```shell
bucketool alias set "minio" -hostname "http://localhost/" -port 9001 -k "minioadmin" -s "minioadmin"

#retour
##Alias has been saved
```

> Une erreur peut survenir lors de l'insertion de l'Alias car un première appel est fait à l'adresse souhaité. Cet appel vous permez d'agir en conséquence si votre alias est défectueux.

Par exemple :

```shell
bucketool alias set "this use cursor" -hostname "http://other.butNotExist/" -port 3000  -k "minioadmin" -s "minioadmin" -c
```

Vous recevrez bien une réponse de type : `Alias has been saved and set as current`
Suivit d'un message d'erreur : `Error while connection:`
Ainsi que du détail de l'erreur : Probablement `Get "http://other.butNotExist:3000/": dial tcp: lookup other.butNotExist: no such host`

Il vous suffira pour cela de rappeller la commande `set` sur le même alias pour le mettre à jour.

#### Alias List

A présent nous pouvons listé les alias que nous avons créer et constaté que le dernier, bien que défectueux est présent.

La command `alias list` comprend très peut de flag ests ont tous optionnels :

- `-detail, -d`_(Optionnel)_ : Permet d'avoir des informations sur l'adresse contactée. Cette dernière ne retourne pas la clef de l'utilisateur ou la clef secret
- `-filter, -f` _(Optionnel)_ : Vous sera utile pour facilité votre recherche.Cette dernière filtreras la liste selon les caractère clef présent. Ce filtre s'applique sur l'alias ainsi que sur l'host.

Afin de facilité la gestion est la lecture des alias, l'alias actuellemnt selectionné sera représenté avec un curseur devant :

##### Exemple d'utilisation de Alias list

```shell
    bucketool alias ls -d -filter "local"

    # Retour
    # Liste was filtered by :  local
    #  - minio (http://localhost:9001/)
```

ou

```shell
bucketool alias ls

# Retour
#   - minio
#   ► this use cursor
```

#### Alias Current

Alias current est une command qui vous permtra deux choses :

- Sans Flag, cette dernière commande retourne l'alias utilisé actuellement
- Avec le flag `-switch`,`-s` suivit du nom d'un autre alias. Ce dernier changera l'adresse du curseur en conséquence.

##### Exemple de l'usage de Alias current

```shell
bucketool alias current

#Retour
# Alias used :  this use cursor
```

Commande avec le flag `-switch`

```shell
bucketool alias current -switch "minio"

#Retour
# Current Alias has been unset : this use cursor
# Switch Alias to minio
```
