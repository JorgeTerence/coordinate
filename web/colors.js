const primary = document.body.dataset.colorPrimary;
const alt = document.body.dataset.colorAlt;

const root = document.querySelector(":root");

const getProp = (el, prop) => getComputedStyle(el).getPropertyValue(prop);
const setProp = (el, prop, v) => el.style.setProperty(prop, v)

document.addEventListener("DOMContentLoaded", () => {
  const clrPrimary = getProp(root, `--${primary}-6`);
  const clrPrimaryDimm = getProp(root, `--${primary}-9`);
  const clrAlt = getProp(root, `--${alt}-5`);

  console.log({ primary, alt, clrPrimary, clrPrimaryDimm, clrAlt });

  if (primary) {
    setProp(root, "--active", clrPrimary);
    setProp(root, "--active-dimm", clrPrimaryDimm);
  }

  if (alt !== undefined) setProp(root, "--alt", clrAlt);
});
