# EN
Simple Discord economy bot

## How to compile?
1. **Clone the repository:**
```sh  
git clone https://github.com/PixelTeam-xyz/PixelEconomy.git  
cd PixelEconomy/  
```  

2. **Check dependencies:**
```sh  
go mod tidy  
```  

3. **Run the build process:**
```sh  
chmod +x ./build.sh  
./build.sh  
```  

4. Done! The binary should be located in the `bin/` directory.

## How to run (for the first time)?
1. **Set up the configuration (optional):**  
   Generate the default configuration:
```sh  
bin/PixelEconomy.elf --initConfig  
```  

If you want, you can now customize it to your needs. More information about this can be found in [CONFIG.md](https://github.com/PixelTeam-xyz/PixelEconomy/blob/main/CONFIG.md).

2. **Provide the bot token:**  
   Create a file named `token.txt` and paste your Discord bot token into it.

3. **Run the bot:**
```sh  
bin/PixelEconomy.elf  
```  

4. Done! The bot should now be running. :)

---

# PL
Prosty bot discord ekonomi

## Jak skompilować?
1. **Sklonuj repozytorium:**
```sh
git clone https://github.com/PixelTeam-xyz/PixelEconomy.git
cd PixelEconomy/
```

2. **Sprawdź zależności:**
```sh
go mod tidy
```

3. **Uruchom kompilacje:**
```sh
chmod +x ./build.sh
./build.sh
```

4. Gotowe! binarka powinna się znaleźć w katalogu bin/

## Jak uruchomić (po raz pierwszy)?
1. **Ustaw konfiguracje (opcjonalnie)**:
   Wygeneruj domyślną konfigurację:
```sh
bin/PixelEconomy.elf --initConfig
```

Jeśli chcesz, możesz ją teraz dostosować do swoich potrzeb, więcej informacji o tym znajdziesz w [CONFIG.md](https://github.com/PixelTeam-xyz/PixelEconomy/blob/main/CONFIG.md)

2. **Podaj token bota:**
   Stwórz plik `token.txt` i wklej do niego swój token bota discord

3. Uruchom bota:
```sh
bin/PixelEconomy.elf
```

4. Gotowe! bot powinien się uruchomić (: