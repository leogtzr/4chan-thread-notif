# 4chan-thread-notif
Little go script that will notify me when there is a new mention to my post in a thread in a 4chan board :)

The script will send you an email to the specified email address.

### Requirements

- To be able to send emails you need a [SENDGRID api key](https://sendgrid.com)
- The program also uses an environmental variable named **EMAIL_TO_4CHAN**

#### How to run it:
1) Passing the environmental variable to the program
```bash
env EMAIL_TO_4CHAN=email@emaildomain.com ./4chan-thread-notifs -board lit -id 15897101 -post 15899180
```
2)
```bash
EMAIL_TO_4CHAN=email@emaildomain.com ./4chan-thread-notifs -board lit -id 15897101 -post 15899180
```
3) Declare your environmental variable in your ~/.bashrc
```bash
EMAIL_TO_4CHAN=email@emaildomain.com
```
source the ~/.bashrc again and run the program:
```bash
./4chan-thread-notifs -board lit -id 15897101 -post 15899180
```

### Build

```
make
```
or:
```
go build -o "4chan-thread-notifs"
```

```bash
./4chan-thread-notifs -board lit -id 15897101 -post 15899180
```
