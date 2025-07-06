### TITLE SCREEN ART ###

title = r'''
  _______ _            _____            _        _   
 |__   __| |          |  __ \          | |      | |  
    | |  | |__   ___  | |__) |___   ___| | _____| |_ 
    | |  | '_ \ / _ \ |  _  // _ \ / __| |/ / _ \ __|
    | |  | | | |  __/ | | \ \ (_) | (__|   <  __/ |_ 
    |_|  |_| |_|\___| |_|  \_\___/ \___|_|\_\___|\__|
               / ____|                               
              | |  __  __ _ _ __ ___   ___           
              | | |_ |/ _` | '_ ` _ \ / _ \          
              | |__| | (_| | | | | | |  __/          
               \_____|\__,_|_| |_| |_|\___|'''

aboutText = r'''
/------------------------------------------------\
| The rocket game is a simple game.              |
| You must succsessfully land on the moon.       |
| But you also need to avoid the atriods and get |
| powerups when you can.                         |
| Along the way you learn many things, make many |
| stories and have lots of fun.                  |
\------------------|  /--------------------------/
     BOB THE 3rd   | /
         _..._     |/
      .'     '.      _
     /    .-""-\   _/ \
   .-|   /:.   |  |   |
   |  \  |:.   /.-'-./
   | .-'-;:__.'    =/
   .'=  *=|NASA _.='
  /   _.  |    ;
 ;-.-'|    \   |
/   | \    _\  _\
\__/'._;.  ==' ==\
         \    \   |
         /    /   /
         /-._/-._/
         \   `\  \
          `-._/._/'''

playSelectedText = r'''
         \\  888b. 8                        
    ------\\ 8  .8 8 .d88 Yb  dP            
    ------// 8wwP' 8 8  8  YbdP             
         //  8     8 `Y88   dP              
             8             dP               
                db    8                 w   
               dPYb   88b. .d8b. 8   8 w8ww 
              dPwwYb  8  8 8' .8 8b d8  8   
             dP    Yb 88P' `Y8P' `Y8P8  Y8P'''

aboutSelectedText = r'''
             888b. 8                        
             8  .8 8 .d88 Yb  dP            
             8wwP' 8 8  8  YbdP             
             8     8 `Y88   dP              
             8             dP               
         \\     db    8                 w   
    ------\\   dPYb   88b. .d8b. 8   8 w8ww 
    ------//  dPwwYb  8  8 8' .8 8b d8  8   
         //  dP    Yb 88P' `Y8P' `Y8P8  Y8P'''

### GAME ART ###

rocks = [
    [
        r" _____",
        r"/     \ ",
        r"|    _/",
        r" \__/",
    ],
    [
        " rrrrr",
        "rrrrrrr",
        "rrrrrrr",
        " rrrr",
    ],
    [
        r"  ___    ___",
        r" /   \__/   \ ",
        r"|        ___/",
        r" \_____/",
    ],
    [
        "  rrr    rrr",
        " rrrrrrrrrrrr",
        "rrrrrrrrrrrrr",
        " rrrrrrr",
    ],
    [
        r" ___",
        r"/   \ ",
        r"\___/",
    ],
    [
        " rrr",
        "rrrrr",
        "rrrrr",
    ],
]

powerups = {
    "health": [
        [
            r"    _",
            r"   | |",
            r"/-------\ ",
            r"|   +   |",
            r"\-------/",
        ],
        [
            "    h",
            "   hhh",
            "hhhhhhhhh",
            "hhhhhhhhh",
            "hhhhhhhhh",
        ],
    ],
    "invinsible": [
        [
            "  /\  ___   /\ ",
            "  \/ /  /   \/",
            "    /  /  /\ ",
            "   /  /   \/",
        ],
        [
            "  ii  iii   ii",
            "  ii iiii   ii",
            "    iiii  ii",
            "   iiii   ii",
        ],
    ],
}

rocketAsciis = {
    "normal": [
        [
            r"       ^  ",
            r"      / \  ",
            r"     / O \  ",
            r"     |   |  ",
            r"    /|   |\  ",
            r"   / |   | \  ",
            r"  /__|   |__\  ",
            r"     |   |     ",
            r"     |   |  ",
            r"    /_____\ ",
        ],
        [
            "       p ",
            "      p p ",
            "     p   p ",
            "     p   p ",
            "    p     p ",
            "   p       p ",
            "  p         p ",
            "     p   p ",
            "     p   p ",
            "    p     p ",
        ],
    ],
    "invinsible": [
        [
            r"       ^ ",
            r"      / \ ",
            r"     /   \ ",
            r"    / ____\ ",
            r"    | \  /| ",
            r"   /|  \/ |\ ",
            r"  / | DEF | \ ",
            r" /__|     |__\ ",
            r"    |     |    ",
            r"   /_______\ ",
        ],
        [
            "       p",
            "      p p",
            "     p   p",
            "    p     p",
            "    p     p",
            "   p       p",
            "  p         p",
            " p           p",
            "    p     p ",
            "   p       p",
        ],
    ],
}
