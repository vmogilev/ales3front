package ales3front;

import org.openqa.selenium.By;
import org.openqa.selenium.WebDriver;
import org.openqa.selenium.WebElement;
import org.openqa.selenium.htmlunit.HtmlUnitDriver;
import static org.junit.Assert.assertTrue;

import org.junit.Test;

public class AppTest {
	private static final String BASEURL = "http://dev.alcalcs.com";

	@Test
    public void testNo() {
        WebDriver driver = new HtmlUnitDriver();
        // test that NO token says "Token validation failed"
        driver.get(BASEURL + "/cdn/uploads/release/CODEGUARDIAN/CODEGUARDIAN_6-7-1_R21.txt.gz?t=NUcYecP0t9Af%2BU0rRvfWh5m3cOVYRmC2c1j9f4YSeQ%3D%3D");
        WebElement element = driver.findElement(By.className("bg-danger"));
        assertTrue(element.getText().contains("Token validation failed"));
        driver.quit();
    }

	@Test
    public void testYes() {
        WebDriver driver = new HtmlUnitDriver();
        // test that YES token says "Start automatic download"
        driver.get(BASEURL + "/cdn/uploads/release/CODEGUARDIAN/CODEGUARDIAN_6-7-1_R21.txt.gz?t=YzJjKhsjV8z0pqacXTk8tc5pYMT2MgLdnHaZPF29fA%3D%3D");
        WebElement element = driver.findElement(By.className("btn-primary"));
        assertTrue(element.getText().contains("Start automatic download"));
        driver.quit();
    }

	@Test
    public void testJunk() {
        WebDriver driver = new HtmlUnitDriver();
        // test that Non-Existent token says "Token validation failed"
        driver.get(BASEURL + "/cdn/uploads/release/CODEGUARDIAN/CODEGUARDIAN_6-7-1_R21.txt.gz?t=Junk");
        WebElement element = driver.findElement(By.className("bg-danger"));
        assertTrue(element.getText().contains("Token validation failed"));
        driver.quit();
    }
	
}
