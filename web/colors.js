const primary = document.body.dataset.colorPrimary;
const alt = document.body.dataset.colorAlt;

const root = document.querySelector(":root");

const getProp = (prop) => getComputedStyle(root).getPropertyValue(prop);
const setProp = (prop, v) => root.style.setProperty(prop, v);

document.addEventListener("DOMContentLoaded", () => {
  const clrPrimary = getProp(`--${primary}-6`);
  const clrPrimaryDimm = getProp(`--${primary}-9`);
  const clrAlt = getProp(`--${alt}-5`);
  
  if (primary) {
    setProp("--active", clrPrimary);
    setProp("--active-dimm", clrPrimaryDimm);
  }

  if (alt !== undefined) {
    setProp("--alt", clrAlt);
  }
});
