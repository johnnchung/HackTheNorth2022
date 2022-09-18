import cm from "../images/cookie.gif"
import dt from "../images/trump.gif"
import bs from "../images/shapiro.gif"
import pg from "../images/penguin.gif"
import gf from "../images/giraffe.gif"
import logo from "../images/snippit-logo.png"
import "./images.css"

export function CookieMonster() {
    return (
        <div class="cm">
            <img src={cm}></img>
        </div>
    )
}

export function DonaldTrump() {
    return (
        <div class="dt">
            <img src={dt}></img>
        </div>
    )
}

export function BenShapiro() {
    return (
        <div class="bs">
            <img src={bs}></img>
        </div>
    )
}

export function Penguin() {
    return (
        <div class="pg">
        <img src={pg}></img>
    </div>
    )
}

export function Giraffe() {
    return (
        <div class="gf">
        <img src={gf}></img>
    </div>
    )
}

export function Logo() {
    return (
        <div class="logo">
        <img src={logo}></img>
    </div>
    )
}