# RemGit

RemGit is a GitCli that lets you search and make repos , search users , clone repos all from within your terminal 

!["Homepage"](images/image1.png)
!["Repository search"](images/image2.png)
!["Repo list after search"](images/image3.png)
!["Repo display"](images/image4.png)
!["User search"](images/image5.png)
!["User list display"](images/image6.png)
!["User repo display"](images/image7.png)


## Installation

***Using Install Script***
```bash
git clone https://github.com/chirag-diwan/RemGit.git
cd RemGit
chmod +x install.sh && ./install.sh

```


***Manual Building***
```bash
git clone https://github.com/chirag-diwan/RemGit.git

cd RemGit && mkdir build && cd build

go build ..

mkdir -p ~/.local/bin

mv RemGit ~/.local/bin

echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc ## or in ~/.zshrc depending on your setup

source ~/.bashrc
```



## Using RemGit

Typing RemGit in your terminal after installation will display the homepage for RemGit,
type `s` at the homepage to enter search mode.

In search mode you can hit `tab` to change search mode between user and repo. Hit `Enter` to start search .

Typing will write text in the searchbar Hitting `Enter` will search for the text in the respective domain (User vs Repository)

After the search result appears you are put into`Normal Mode` where you could navigate the list using `j` or `down` for down and `k` or `up` for up.

Hiting `Enter` on a search result will open details for the selected item

Hiting `c` on a repository item will clone it into the directory you are in.

Hiting `backspace` on details page will navigate you back 

## Configuration

The config lives in `~/.remgit.conf` file on you system , below is a overview of the configuration that is supported.

```conf
//Your Github Personal Access token needed for creating repos

PAT=<github personal access token>

//Toggle Homepage display at start
Showhome = true 

//Change the color display
Subtle    = "#45475a" 
Highlight = "#cba6f7"
Text      = "#cdd6f4"
Warning   = "#f38ba8"
Special   = "#a6e3a1"
```
