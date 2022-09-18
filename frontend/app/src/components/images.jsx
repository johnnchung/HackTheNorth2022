import cm from "../images/cookie.gif";
import dt from "../images/trump.gif";
import bs from "../images/shapiro.gif";
import pg from "../images/penguin.gif";
import gf from "../images/giraffe.gif";
import logo from "../images/snippit-logo.png";
import "./images.css";

export function CookieMonster() {
  return (
    <div className="cm">
      <img src={cm} alt="Cookie Monster"></img>
    </div>
  );
}

export function DonaldTrump() {
  return (
    <div className="dt">
      <img src={dt} alt="Donald Trump"></img>
    </div>
  );
}

export function BenShapiro() {
  return (
    <div className="bs">
      <img src={bs} alt="Ben Shapiro"></img>
    </div>
  );
}

export function Penguin() {
  return (
    <div className="pg">
      <img src={pg} alt="Penguin"></img>
    </div>
  );
}

export function Giraffe() {
  return (
    <div className="gf">
      <img src={gf} alt="Giraffe"></img>
    </div>
  );
}

export function Logo() {
  return (
    <div className="logo">
      <img src={logo} alt="logo"></img>
    </div>
  );
}
