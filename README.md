How to run playbook:

    $ ansible-playbook -i 192.168.0.104, gorshock.yml

Build bot:

    $ cd bot && GOOS=linux GOARCH=arm GOARM=6 go build -v .

Deploy bot:

    $ scp bot/bot pi@192.168.19.50:.
