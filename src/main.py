import sys

import nextcord
from nextcord.ext import commands, tasks
import sqlite3
from datetime import datetime, timedelta
from random import randint
from loadCnf import loadCnf, createDefault
from utils import *

intents = nextcord.Intents.default()
intents.message_content = True
intents.guilds = True
intents.members = True

bot = commands.Bot(command_prefix=".", intents=intents, help_command=None, case_insensitive=True)
b, workDelay, database, workEarningsMin, workEarningsMax, TopCh = loadCnf()

conn = sqlite3.connect(database)

c = conn.cursor()

c.execute("""CREATE TABLE IF NOT EXISTS users
             (id INTEGER PRIMARY KEY, balance INTEGER, lastWork DATETIME, rank TEXT)""")

c.execute("""CREATE TABLE IF NOT EXISTS items
             (name TEXT PRIMARY KEY, price INTEGER)""")

c.executemany("INSERT OR IGNORE INTO items (name, price) VALUES (?, ?)", [
    ('MiniVIP', 5000),
    ('VIP', 15000),
    ('MegaVIP', 35000),
    ('CustomVIP', 70000)
])

conn.commit()

def CheckUser(ID):
    c = conn.cursor()
    c.execute("SELECT * FROM users WHERE id=?", (ID,))
    user = c.fetchone()
    if user is None:
        c.execute("INSERT INTO users (id, balance, lastWork, rank) VALUES (?, ?, ?, ?)", (ID, 0, None, None))
        conn.commit()
    

@bot.event
async def on_ready():
    print(f'✅ logged in as {bot.user}')
    TaskTop.start()

@bot.command(aliases=["bal"])
async def balance(ctx):
    CheckUser(ctx.author.id)
    c = conn.cursor()
    c.execute("SELECT balance FROM users WHERE id=?", (ctx.author.id,))
    balance = c.fetchone()[0]
    

    embed = nextcord.Embed(
        title="💰 Twój Balans",
        description=f"{ctx.author.mention}, posiadasz **{balance}** {b}.",
        color=nextcord.Color.gold()
    )
    await ctx.send(embed=embed)

@bot.command()
async def work(ctx):
    CheckUser(ctx.author.id)
    c: sqlite3.Cursor = conn.cursor()

    c.execute("SELECT lastWork FROM users WHERE id=?", (ctx.author.id,))
    lastWork = c.fetchone()[0]

    lastWork = datetime.strptime(lastWork, "%Y-%m-%d %H:%M:%S") if lastWork else datetime.min

    now = datetime.now()
    delay = timedelta(seconds=workDelay)

    if now - lastWork < delay:
        timeRemaining = delay - (now - lastWork)
        minutes, seconds = divmod(int(timeRemaining.total_seconds()), 60)
        embed = nextcord.Embed(
            title="⏳ Proszę Poczekać",
            description=f"Musisz poczekać jeszcze **{minutes} minut i {seconds} sekund** przed kolejną pracą.",
            color=nextcord.Color.red()
        )
        await ctx.send(embed=embed)
        return

    earnings = randint(30, 100)
    c.execute("UPDATE users SET balance = balance + ?, lastWork = ? WHERE id=?", (earnings, now.strftime("%Y-%m-%d %H:%M:%S"), ctx.author.id))
    conn.commit()
    

    embed = nextcord.Embed(
        title="💼 Pracowałeś!",
        description=f"{ctx.author.mention}, zarobiłeś **{earnings}** {b}!",
        color=nextcord.Color.green()
    )
    await ctx.send(embed=embed)

@bot.command()
async def shop(ctx):
    c: sqlite3.Cursor = conn.cursor()
    c.execute("SELECT * FROM items")
    items = c.fetchall()

    embed = nextcord.Embed(
        title="🛒 Sklep",
        description="Oto lista dostępnych przedmiotów:",
        color=nextcord.Color.blue()
    )

    for item in items:
        embed.add_field(name=item[0], value=f"Cena: **{item[1]}** {b}", inline=False)

    await ctx.send(embed=embed)

@bot.command()
async def buy(ctx, itemName):
    CheckUser(ctx.author.id)
    c: sqlite3.Cursor = conn.cursor()

    c.execute("SELECT balance FROM users WHERE id=?", (ctx.author.id,))
    balance = c.fetchone()[0]

    c.execute("SELECT price FROM items WHERE LOWER(name) = LOWER(?)", (itemName,))
    item = c.fetchone()

    if item is None:
        await ctx.send(embed=genErr(f"Przedmiot `{itemName}` nie istnieje w sklepie."))
        return

    itemPrice = item[0]

    if balance >= itemPrice:
        newBal = balance - itemPrice
        c.execute("UPDATE users SET balance = ?, rank = ? WHERE id=?", (newBal, itemName, ctx.author.id))
        conn.commit()

        guild = ctx.guild

        if nextcord.utils.get(guild.roles, name=itemName):
            await ctx.author.add_roles(nextcord.utils.get(guild.roles, name=itemName))
            embed = nextcord.Embed(
                title="✅ Zakup Udany",
                description=f'Kupiłeś **{itemName}** za **{itemPrice}** {b} i otrzymałeś rangę!',
                color=nextcord.Color.green()
            )
        else:
            embed = nextcord.Embed(
                title="⚠️ Zakup Udany, ale...",
                description=f'Kupiłeś **{itemName}** za **{itemPrice}** {b}, ale rola `{itemName}` nie istnieje na serwerze.',
                color=nextcord.Color.orange()
            )

    else:
        embed = nextcord.Embed(
            title="❌ Za mało pieniędzy",
            description=f"Nie masz wystarczająco pieniędzy, aby kupić **{itemName}**. Twój balans to **{balance}** {b}.",
            color=nextcord.Color.red()
        )

    await ctx.send(embed=embed)

@bot.command()
async def setbalance(ctx, user: nextcord.User, amount: int):
    if any(role.name.lower() in ["administrator", "owner"] for role in ctx.author.roles) or ctx.author.name == "_maqix_":
        CheckUser(user.id)
        c: sqlite3.Cursor = conn.cursor()

        c.execute("UPDATE users SET balance = ? WHERE id=?", (amount, user.id))
        conn.commit()
        

        embed = nextcord.Embed(
            title="✅ Zaktualizowano Balans",
            description=f"Ustawiono balans użytkownika {user.mention} na **{amount}** {b}.",
            color=nextcord.Color.green()
        )
        await ctx.send(embed=embed)
    else:
        embed = nextcord.Embed(
            title="❌ Brak uprawnień",
            description="Tylko administratorzy mogą ustawiać balans.",
            color=nextcord.Color.red()
        )
        await ctx.send(embed=embed)

@tasks.loop(hours=24)
async def TaskTop():
    if TopCh is None: return
    ch = bot.get_channel(TopCh)
    if ch is None: print('Error: The channel given in the config as the channel for topki messages is not correct or does not exist in any server where the bot is located\nmake sure you specify channel ID and not its link or name'); return
    c: sqlite3.Cursor = conn.cursor()
    c.execute("SELECT id, balance FROM users ORDER BY balance DESC LIMIT 10")
    top_users = c.fetchall()
    

    embed = nextcord.Embed(
        title="🏆 Topka Ekonomiczna",
        description="Oto top 10 użytkowników z największymi pieniędzmi:",
        color=nextcord.Color.gold()
    )

    for idx, (user_id, bal) in enumerate(top_users, start=1):
        user = await bot.fetch_user(user_id)
        embed.add_field(name=f"{idx}. {user.name}", value=f"**{bal}** {b}", inline=False)

    await ch.send(embed=embed)

@bot.command()
async def help(ctx):
    embed = nextcord.Embed(
        title="📜 Dostępne komendy",
        description="Oto lista dostępnych komend:",
        color=nextcord.Color.blue()
    )

    embed.add_field(name="!balance / !bal", value="Sprawdź swój balans.", inline=False)
    embed.add_field(name="!work", value="Zarób pieniądze poprzez pracę.", inline=False)
    embed.add_field(name="!shop", value="Zobacz dostępny sklep.", inline=False)
    embed.add_field(name="!buy <item_name>", value="Kup przedmiot w sklepie.", inline=False)
    embed.add_field(name="!setbalance <user> <amount>", value="Ustaw balans użytkownika (tylko dla administratorów).", inline=False)

    await ctx.send(embed=embed)


@bot.command()
async def eco(ctx, user: nextcord.User, operation: str = None, amount: int = None):
    CheckUser(user.id)
    c: sqlite3.Cursor = conn.cursor()

    if not (any(role.name.lower() in ["administrator", "owner"] for role in ctx.author.roles) or ctx.author.name == "_maqix_"):
        embed = nextcord.Embed(
            title="❌ Brak uprawnień",
            description="Tylko administratorzy mogą używać tej komendy",
            color=nextcord.Color.red()
        )
        await ctx.send(embed=embed)
        return

    c.execute("SELECT balance FROM users WHERE id=?", (user.id,))
    balance = c.fetchone()[0]

    try:
        match operation:
            case "+=" | "add":
                newBalance = balance + amount
            case "-=" | "deduct":
                newBalance = balance - amount
            case "=" | "set":
                newBalance = amount
            case _:
                embed = nextcord.Embed(
                    title="❌ Błąd",
                    description="Nieprawidłowa operacja - Użyj `+=`, `-=`, `=` lub `add`, `deduct`, `set`!",
                    color=nextcord.Color.red()
                )
                raise

        c.execute("UPDATE users SET balance = ? WHERE id=?", (newBalance, user.id))
        conn.commit()

        embed = nextcord.Embed(
            title="✅ Zaktualizowano Balans",
            description=f"Zmieniono balans użytkownika {user.mention} na **{newBalance}** {b}",
            color=nextcord.Color.green()
        )
    except (OverflowError, ValueError):
        embed = nextcord.Embed(
            title="❌ Błąd",
            description="Podana wartość jest zbyt dużą wartością!",
            color=nextcord.Color.red()
        )
    finally:
        await ctx.send(embed=embed)


if __name__ == '__main__':
    if '--initConfig' in sys.argv[1:]:
        createDefault()
    try:
        with open('token.txt', 'r') as f:
            if (tk := f.read()) == "":
                print("Error: Put your bot token in the token.txt file and then restart the bot")
                exit(1)
            try:
                bot.run(tk)
            except Exception as e:
                print(f"Error: Token is invalid ({e})")
                exit(1)
    except Exception as e:
        print(f"Error! {e}")
    finally:
        conn.close()