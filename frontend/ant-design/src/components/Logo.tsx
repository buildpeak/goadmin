import logo from "../assets/images/ding.svg";

type LogoProps = {
  width: number;
};

const Logo = ({ width }: LogoProps) => {
  return <img width={width} src={logo} alt="logo" />;
};

export default Logo;
