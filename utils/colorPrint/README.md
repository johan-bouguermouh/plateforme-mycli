# Usage du Package ToolPrint

## Avant Propos

Ce petit package utile offre rapidement la possibilité à qui veux de structurer les logs retournées dans le terminal de manière à amélioré sa lecture. Son approche est seulement de facilité l'encription des conde ANSI dans le terminal. Par consequent, mise à part la dernière fonction proposée dans le package, toutes les autre son appliquer dans un fonction d'impretion fournis pas `fmt`.

## Utilisation

### Fonctions simple de colorisation

Ces fonctions sont utilisé pour permttre d'intégré rapidement un couleur sans trop du diffculter.Elle en comportent pas de casse particulière mais seulement l'application de couleurs.

#### Gloassaire des fonctions de type simple

| Appel     | Couleur | attributs de la fonction |
| --------- | ------- | ------------------------ |
| `BlueP`   | Blue    | **content** _string_     |
| `GrennP`  | Vert    | **content** \_string     |
| `YellowP` | Jaune   | **content** \_string     |
| `RedP`    | Rouge   | **content** \_string     |
| `GreyP`   | Gris    | **content** \_string     |
| `PurpleP` | Violet  | **content** \_string     |

Lors de leurs usage dans une impression vous pouvez très simplement concatainer les différentes fonctions dans un même print

```go
package main

import (
    "fmt"
    cp "api-interface/utils/colorPrint"
)
func main {
     fmt.PrintLn(cp.BlueP("Je commence en Bleu"), cp.GreenP(", je fini en rouge"))
}
```

### Fonction spécifique : ColorPrint

Si vous souhaitez allé plus loin sur l'usage vous pouvez aussi utiliser un fonction qui permet de cummuler le type de casses et les couleurs associer. Cette fonctionnaliter et plus verbeuse et peut évoluée dans le temps mais l'approche est simple :

La méthode comprends trois attributs :

- **color** _string_ : Détermine le code couleur de la font
- **Content** _string_ : Votre contenu
- **Options** \*colorPrint.Options\*\* : Différente options à appliquer

#### Type d'options possible

| Attribut   | Type    | Définition                                                            |
| ---------- | ------- | --------------------------------------------------------------------- |
| Bold       | `bool`  | si true, renvoie la chaine de caractère en gras                       |
| Italic     | `bool`  | si true, renvoie la chaine de caractère en italic                     |
| Underline  | `bool`  | Si, true renvoie la châine de caractère souligné                      |
| Background | `Strin` | Si déclarée, renvoie la couleur présente dans la valeur de l'attribut |

#### Usage

```go
package main

import (
    "fmt"
    cp "api-interface/utils/colorPrint"
)
func main {
     	fmt.Println(colorPrint.ColorPrint("Blue", "Hello, World!", &colorPrint.Options{
        Background: "Blue",
        Bold:       true,
    }))
    fmt.Println(colorPrint.ColorPrint("Green", "Hello, World!", &colorPrint.Options{
        Italic:     true,
        Underline:  true,
        Background: "Green",
    }))
}
```

### Function ReadObject (_beta_)

Vous avez accès aussi a un fonction béta qui permet de lire n'importe qu'elle type d'object tant que ce dernier compren des attributs publics. Cette fonction permet d'avoir un retour plus prorpe dans le terminal concernant un obkect que nous souhaitons consulté.
