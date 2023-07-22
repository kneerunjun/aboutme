Now that I have decided to give vuejs a try , I was trying to get vuecli installed on my machine. - but then that means I have to go thru th  nodejs and npm shit.  I had hiccups installing nodejs and npm on my machine before switching (as expected, never does the js community provide a clean installation), but nevertheless got it running 

```
yay -S nvm  
echo 'source /usr/share/nvm/init-nvm.sh' >> ~/.zshrc

```
> DO NOT FORGET : start your terminal again 

```
nvm -v
nvm install --lts # this is to install node and npm
node -v 
npm -v 

```
Boom ! nvm is installed on your system. You should be able to see the versions of node and npm respectively

```
 npm install -g @vue/cli@4.5.13
```
This shall install vue cli globally on your system. So that does it you are now ready to create a new project