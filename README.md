# plateforme-mycli

## Initialisation du projet

Pour installer la CLI, exécutez simplement :

```
go get github.com/urfave/cli
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

### Flags Généraux

Vous pouvez exploiter l'application de différentes manière. L'une des première approche est de définir un alias qui sera exploité par défaut lors de vos différente commande. Afin de comprendre correctement son usage nous commencerons par cela. Toutes les commandes suivantes seront déterminés par la bonne insertion de vos alias. Assurez-vous que ces derniers soient correct est adaptés à l'usage d'un bucket S3.

Si vous souhaitez cependant altrernée l'usage du CLI sur plusieur bucket vous pouvez spécifier l'alias cible avec l'usage du flags `-alias`.

```shell
bucketool -alias <command> ...
```

Vous pouvez notamment passer le programme en mod debug avec l'alias `-debug`. Ainsi, vous aurez plus d'information sur les commands appeler et le fonctionnement interne de l'application. Lors de l'usage du flag `-debug`, une fonction d'interception de requêtes HTPP vous permetra d'avoir un vision plus détaillées sur ce qu'il ce passe.

```shell
bucketool -debug <command> ...
```

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
- _`-keyname, -k_`**(obligatoire)**: Définit le AWSAccessKeyId
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

#### Commande `alias delete`

La commande `alias delete` permet de supprimer un alias de la liste des alias. Vous pouvez spécifier le nom de l'alias à supprimer ou utiliser des options pour supprimer tous les alias ou ceux qui ne peuvent pas se connecter au serveur.

##### Usage

```sh
bucketool alias delete <name>
```

##### Exemple de l'usage de Alias delete

```sh
bucketool alias delete myalias
```

##### Options

- -a, --all : Supprime tous les alias.
- -sc, --savecurrent : Sauvegarde l'alias actuel lorsque vous utilisez l'option -a ou --all.
- -c, --clean : Supprime tous les alias qui ne peuvent pas se connecter au serveur.

##### Exemple d'utilisation avec options

###### Supprimer tous les alias

```sh
bucketool alias delete --all
```

###### Supprimer tous les alias mais sauvegarder l'alias actuel

```sh
bucketool alias delete --all --sc
```

###### Supprimer tous les alias qui ne peuvent pas se connecter au serveur

```sh
bucketool alias delete --clean
```

##### Confirmation de suppression

Si l'alias actuel est sur le point d'être supprimé, une confirmation sera demandée à l'utilisateur :

```text
The current Alias same as touched by this command will be deleted, do you want to delete ? (y/n) :
```

L'utilisateur doit répondre par y pour confirmer la suppression ou par n pour conserver l'alias.

Si l'alias actuel est supprimé, un message d'avertissement sera affiché :
`WARN | The current Alias has been deleted`

### Usage des commandes Bucket

#### Command Bucket Create

La commande `create bucket` permet de créer un nouveau bucket sur un serveur compatible S3 (comme MinIO). Cette commande vérifie d'abord si le nom du bucket est valide, puis tente de créer le bucket. Si le bucket existe déjà, un message d'erreur approprié est affiché.

##### Arguments de la command Bucket Create

- `<name>` : Le nom du bucket à créer. Cet argument est obligatoire et doit être unique. Le nom doit être en minuscules et ne contenir que des lettres minuscules, des chiffres, des tirets (-) et des points (.).

1. Le nom ne doit pas être vide.
2. Le nom doit comporter au moins 3 caractères.
3. Le nom doit comporter au plus 63 caractères.
4. Le nom doit être en minuscules.
5. Le nom ne doit contenir que des lettres minuscules, des chiffres, des tirets (-) et des points (.).

##### Options de la command Bucket Create

-alias : Nom de l'alias à utiliser. Si vous avez spécifié l'alias actuel, vous pouvez omettre cette option. L'utilisation de cette option est facultative. Si vous l'utilisez, cette option doit être placée après "bucket" et avant "create `<name>`", comme ceci : `-alias <alias> bucket create <name>`.

##### Exemple de la command Bucket Create

```shell
bucketool bucket create mybucket
bucketool -alias "myalias" bucket create "mybucket"
```

#### Commande `bucket list`

La commande `bucket ls` permet de lister tous les buckets sur un serveur compatible S3 (comme MinIO). Cette commande peut également afficher des informations détaillées sur chaque bucket si l'option `--details` est spécifiée.

##### Arguments de la commande bucket list

Aucun argument n'est requis pour cette commande.

##### Options de la commande bucket list

- `-d`, `--details` : Affiche des informations détaillées sur chaque bucket, y compris la région, les ACL (Access Control List), la journalisation et la version.

##### Exemple de la commande bucket liste

Pour lister tous les buckets :

```shell
bucketool bucket ls
```

ou

```shell
bucketool bucket ls -d
```

#### Commande `bucket delete`

La commande `bucket delete` permet de supprimer un bucket sur un serveur compatible S3 (comme MinIO). Cette commande vérifie d'abord si le bucket existe, puis tente de le supprimer. Si le bucket n'existe pas, un message d'erreur approprié est affiché.

##### Arguments de la commande bucket delete

<name> : Le nom du bucket à supprimer. Cet argument est obligatoire et doit correspondre à un bucket existant.

##### Options de la commande bucket delete

`-alias` : Nom de l'alias à utiliser. Si vous avez spécifié l'alias actuel, vous pouvez omettre cette option. L'utilisation de cette option est facultative. Si vous l'utilisez, cette option doit être placée après "bucket" et avant "delete <name>", comme ceci : `-alias <alias>  bucket delete <name>`.

##### Exemples de la commande bucket delete

Pour supprimer un bucket nommé `mybucket` :

```shell
bucketool bucket delete mybucket
```

ou

```shell
bucketool bucket d "mybucket"
```

### Usage des commandes de Bucket Object

#### Commande `list_objects`

La commande list_objects permet de lister tous les objets dans un bucket S3. Elle offre également une option pour afficher des détails supplémentaires sur chaque objet.

##### Options de `list_objects`

- `-b`, `--bucket` _(Requis)_ : Nom du bucket de destination.
- `-d`, `--details` _(Optionnel)_ : Afficher les détails des objets.
- `-alias` _(Optionnel)_ : Spécifier un alias à utiliser. Si vous avez spécifié l'alias actuel, vous pouvez omettre cette option.

##### Exemple d'usage de `list`, `ls`

- Lister tous les objets dans un bucket :

```shell
bucketool ls -b mybucket

```

- Lister tous les objets dans un bucket avec des détails :

```shell
bucketool list -b mybucket -d
```

Utiliser un alias pour spécifier le bucket :

```shell
bucketool -alias "myalias" ls -b "mybucket"
```

#### Commande `copy`, `cp`

La commande `copy` permet de copier un fichier depuis un chemin local et de l'insérer dans un bucket de destination sur un serveur compatible S3.

##### Options de `copy`

- `-d`, `--destination` _(Requis)_ : Nom du bucket de destination.
- `-n`, `--name` _(Optionnel)_ : Nom du fichier dans le bucket. Si non spécifié, le nom du fichier sera le même que celui du fichier copié.
- `-alias` _(Optionnel)_ : Spécifier un alias à utiliser. Si vous avez spécifié l'alias actuel, vous pouvez omettre cette option.

##### Exemple d'usage de `copy`

- Copier un fichier dans un bucket en concervant le nom d'origine

```shell
bucketool cp "/path/to/file" -d "mybucket"
```

- Copier un fichier dans un bucket avec un nom spécifique

```shell
bucketool copy "/path/to/file" -d "mybucket" -n "myfile.txt"
```

#### Commande `download`, `dl`

La commande `download` permet de télécharger un fichier depuis un bucket S3 et de l'insérer dans un chemin local. Elle détecte automatiquement le type MIME du fichier et ajoute l'extension appropriée au fichier téléchargé.

#### Options de `download`

- `-b`, `--bucket` _(Requis)_ : Nom du bucket où se trouve le fichier.
- `-n`, `--name` _(Requis)_ : Nom du fichier dans le bucket.
- `-rename`, `-rn` (_Optionnel)_ : Nom du fichier dans le chemin local. Vous n'avez pas besoin de spécifier l'extension, elle sera ajoutée automatiquement. Si vous ne la spécifiez pas, le nom du fichier sera le même que celui du fichier copié, mais il pourrait être modifié si le type MIME est différent.
- `-alias` _(Optionnel)_ : Spécifier un alias à utiliser. Si vous avez spécifié l'alias actuel, vous pouvez omettre cette option.

#### Exemple d'usage de `download`

- Télécharger un fichier dans un bucket

```shell
bucketool download "/path/to/file" -b mybucket -n myfile.txt
```

- Télécharger un fichier dans un bucket avec un nouveau nom

```shell
 bucketool /path/to/file -b mybucket -n myfile.txt -rename newfile
```

Utiliser un alias pour spécifier le bucket :

```shell
 bucketool download -alias myalias /path/to/file -b mybucket -n myfile.txt -rename newfile
```

#### Commande `del`

La commande `delete` permet de supprimer un objet existant d'un bucket existant. Cette commande nécessite le nom du bucket à supprimer comme argument ainsi que le nom d'lobjet cible. Vous pouvez également utiliser l'option -alias pour spécifier un alias avant la commande `delete`.

#### Options de `del`

- `-bucket`, `-b` _(required)_ : Le nom du bucket de destination. Ce flag est obligatoire.
- `-name`, `-n` _(required)_ : Le nom du fichier dans le bucket. Ce flag est obligatoire.
- `-alias` _(Optionnel)_ : Spécifier un alias à utiliser. Si vous avez spécifié l'alias actuel, vous pouvez omettre cette option.

#### Exemple d'usage de `del`

- Avec alias selectionné

```shell
bucketool del --bucket "mybucket" --name "myfile.txt"
```

- En ajoutant un alias manuellement

```shell
bucketool -alias "myAlias" del -b "mybucket" -n "myfile.txt"
```
