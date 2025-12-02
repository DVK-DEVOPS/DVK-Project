import { expect, test } from "@playwright/test";

test("homepage loads", async ({ page }) => {
  await page.goto("/");
  await expect(page.getByRole("textbox", { name: "Search..." })).toBeVisible();
});

test("search works", async ({ page }) => {
  await page.goto("/");
  await page.getByRole("textbox", { name: "Search..." }).fill("java");
  await page.getByRole("button", { name: "Search" }).click();
  await expect(
    page.getByText("Java Streams Examples v2 progurl2.dev-v2")
  ).toBeVisible();
});

test("navigation works", async ({ page }) => {
  await page.goto("/");
  await page.getByRole("link", { name: "Register" }).click();
  await expect(page).toHaveURL(/register/);
  await page.getByRole("link", { name: "Already have an account?" }).click();
  await expect(page).toHaveURL(/login/);
});

test("register validation", async ({ page }) => {
  const uid = Date.now().toString();
  await page.goto("/register");
  await page.locator('input[name="username"]').fill(`test_${uid}`);
  await page.locator('input[name="email"]').fill(`test_${uid}@test.com`);
  await page.locator('input[name="password"]').fill("test33");
  await page.locator('input[name="password2"]').fill("test33");
  await page.getByRole("button", { name: "Sign up" }).click();
  await page.waitForURL("/");
  await expect(page.getByRole("textbox", { name: "Search..." })).toBeVisible();
});

test("login flow", async ({ page }) => {
  const uid = Date.now().toString();
  await page.goto("/register");
  await page.locator('input[name="username"]').fill(`u${uid}`);
  await page.locator('input[name="email"]').fill(`u${uid}@t.com`);
  await page.locator('input[name="password"]').fill("test33");
  await page.locator('input[name="password2"]').fill("test33");
  await page.getByRole("button", { name: "Sign up" }).click();
  await page.waitForURL("/");

  await page.goto("/login");
  await page.locator('input[name="username"]').fill(`u${uid}`);
  await page.locator('input[name="password"]').fill("test33");
  await page.getByRole("button", { name: "Log in" }).click();

  await expect(page.getByRole("textbox", { name: "Search..." })).toBeVisible();
});
