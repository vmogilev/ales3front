package ales3front;

import org.openqa.selenium.By;
import org.openqa.selenium.WebDriver;
import org.openqa.selenium.WebElement;
import org.openqa.selenium.htmlunit.HtmlUnitDriver;
//import org.openqa.selenium.firefox.FirefoxDriver;
import static org.junit.Assert.assertTrue;

import org.junit.Test;

public class AppTest {
	// And Example of ALE3FRONT_TEST_HTTP_HOST is "http://localhost:8080";
	private static final String BASEURL = System.getenv("ALE3FRONT_TEST_HTTP_HOST");
	private static final String FAIL_MESSAGE = "Token validation failed";
	private static final String DOWNLOAD_MESSAGE = "Start download";
	private static final String DEADLINE_MESSAGE = "context deadline exceeded";
	private static final String DOWNLOAD_PATH = "/cdn/uploads/release/CODEGUARDIAN/CODEGUARDIAN_6-7-1_R21.txt.gz";
	private static final String NO_TOKEN = "NUcYecP0t9Af%2BU0rRvfWh5m3cOVYRmC2c1j9f4YSeQ%3D%3D";
	private static final String YES_TOKEN = "YzJjKhsjV8z0pqacXTk8tc5pYMT2MgLdnHaZPF29fA%3D%3D";

	@Test
    public void testNo() {
        WebDriver driver = new HtmlUnitDriver();
        // test that NO token says "Token validation failed"
        driver.get(BASEURL + DOWNLOAD_PATH + "?t=" + NO_TOKEN);
        WebElement element = driver.findElement(By.id("dl-error"));
        assertTrue(element.getText().contains(FAIL_MESSAGE));
        driver.quit();
    }

	@Test
    public void testYes() {
        WebDriver driver = new HtmlUnitDriver();
        // test that YES token says "Start automatic download"
        driver.get(BASEURL + DOWNLOAD_PATH + "?t=" + YES_TOKEN);
        WebElement element = driver.findElement(By.id("dl-button"));
        assertTrue(element.getText().contains(DOWNLOAD_MESSAGE));
        driver.quit();
    }

	@Test
    public void testJunk() {
        WebDriver driver = new HtmlUnitDriver();
        // test that Non-Existent token says "Token validation failed"
        driver.get(BASEURL + DOWNLOAD_PATH + "?t=Junk");
        WebElement element = driver.findElement(By.id("dl-error"));
        assertTrue(element.getText().contains(FAIL_MESSAGE));
        driver.quit();
    }
	
	@Test
    public void testTimeout() {
        WebDriver driver = new HtmlUnitDriver();
        // test that YES token with short Auth Timeout says "context deadline exceeded"
        driver.get(BASEURL + DOWNLOAD_PATH + "?o=10&t=" + YES_TOKEN);
        WebElement element = driver.findElement(By.id("dl-error"));
        assertTrue(element.getText().contains(DEADLINE_MESSAGE));
        driver.quit();
    }

}
