# retry

`retry` a command line until it worked (exit code=0), it printed something or it exceeds a timeout.

# Usage

```sh
$ retry -h

retry [options] [cmd]

Usage of retry:
  -ok-text string
    	a text to lookup for inthe stdout/stderr that confirms did run correctly.
  -retry-interval duration
    	Retry interval duration (default 1s)
  -tail
    	Apply interval at tail or head
  -timeout duration
    	maximum timeout duration (default 1m0s)
```

# Example

```sh
$ retry ls -alh
2018/06/10 15:01:31 => ls  -alh
total 20K
drwxrwxr-x    3 mh-cbon mh-cbon 4,0K 10 juin  14:52 .
drwxrwxr-x. 116 mh-cbon mh-cbon 4,0K 10 juin  14:52 ..
drwxrwxr-x    7 mh-cbon mh-cbon 4,0K 10 juin  15:01 .git
-rw-rw-r--    1 mh-cbon mh-cbon 3,2K 10 juin  15:01 main.go
-rw-rw-r--    1 mh-cbon mh-cbon  108 10 juin  14:55 README.md

$ retry -timeout 5s wget nop
2018/06/10 15:01:14 => wget  nop
--2018-06-10 15:01:14--  http://nop/
Résolution de nop (nop)… échec : Name or service not known.
wget : impossible de résoudre l’adresse de l’hôte « nop »
2018/06/10 15:01:15 => wget  nop
--2018-06-10 15:01:15--  http://nop/
Résolution de nop (nop)… échec : Name or service not known.
wget : impossible de résoudre l’adresse de l’hôte « nop »
2018/06/10 15:01:16 => wget  nop
--2018-06-10 15:01:16--  http://nop/
Résolution de nop (nop)… échec : Name or service not known.
wget : impossible de résoudre l’adresse de l’hôte « nop »
2018/06/10 15:01:17 => wget  nop
--2018-06-10 15:01:17--  http://nop/
Résolution de nop (nop)… échec : Name or service not known.
wget : impossible de résoudre l’adresse de l’hôte « nop »
2018/06/10 15:01:18 => wget  nop
2018/06/10 15:01:18 failed to execute "wget" [nop], err=context deadline exceeded

$ retry -ok-text html wget -O /dev/null http://google.com
2018/06/10 15:10:03 => wget  -O /dev/null http://google.com
--2018-06-10 15:10:03--  http://google.com/
Résolution de google.com (google.com)… 216.58.207.142
Connexion à google.com (google.com)|216.58.207.142|:80… connecté.
requête HTTP transmise, en attente de la réponse… 301 Moved Permanently
Emplacement : http://www.google.com/ [suivant]
--2018-06-10 15:10:03--  http://www.google.com/
Résolution de www.google.com (www.google.com)… 216.58.207.132
Connexion à www.google.com (www.google.com)|216.58.207.132|:80… connecté.
requête HTTP transmise, en attente de la réponse… 200 OK
Taille : non indiqué [text/html]
Sauvegarde en : « /dev/null »

     0K .......... ..                                           717K=0,02s

2018-06-10 15:10:03 (717 KB/s) - « /dev/null » sauvegardé [13120]
```

# Install

```sh
go get github.com/mh-cbon/retry
```
