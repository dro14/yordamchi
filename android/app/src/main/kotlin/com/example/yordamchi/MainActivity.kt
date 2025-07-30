package uz.chuqurtech.yordamchi

import android.content.ComponentName
import android.content.pm.PackageManager
import io.flutter.embedding.android.FlutterActivity
import io.flutter.embedding.engine.FlutterEngine
import io.flutter.plugin.common.MethodChannel

class MainActivity : FlutterActivity() {
    private val CHANNEL = "uz.chuqurtech.yordamchi/dynamic_icon"

    override fun configureFlutterEngine(flutterEngine: FlutterEngine) {
        super.configureFlutterEngine(flutterEngine)
        
        MethodChannel(flutterEngine.dartExecutor.binaryMessenger, CHANNEL).setMethodCallHandler { call, result ->
            when (call.method) {
                "setAlternateIconName" -> {
                    val newIcon = call.argument<String>("newIcon")
                    val oldIcon = call.argument<String>("oldIcon")
                    if (newIcon != null && oldIcon != null) {
                        try {
                            setAlternateIcon(newIcon, oldIcon)
                            result.success(null)
                        } catch (e: Exception) {
                            result.error("ICON_ERROR", "Failed to set icon: ${e.message}", null)
                        }
                    } else {
                        result.error("INVALID_ARGUMENT", "newIcon and oldIcon are required", null)
                    }
                }
                else -> {
                    result.notImplemented()
                }
            }
        }
    }

    private fun setAlternateIcon(newIcon: String, oldIcon: String) {
        val packageManager = packageManager
        val packageName = packageName

        val iconAliases = listOf(
            "pink", "red", "deepOrange", "orange", "amber", "yellow", "lime", 
            "lightGreen", "green", "teal", "cyan", "lightBlue", "blue", 
            "indigo", "purple", "deepPurple", "blueGrey", "brown", "grey"
        )

        if (oldIcon == "default") {
            val mainActivityName = ComponentName(packageName, "$packageName.MainActivity")
            packageManager.setComponentEnabledSetting(
                mainActivityName,
                PackageManager.COMPONENT_ENABLED_STATE_DISABLED,
                PackageManager.DONT_KILL_APP
            )
        } else if (iconAliases.contains(oldIcon)) {
            val oldComponentName = ComponentName(packageName, "$packageName.$oldIcon")
            packageManager.setComponentEnabledSetting(
                oldComponentName,
                PackageManager.COMPONENT_ENABLED_STATE_DISABLED,
                PackageManager.DONT_KILL_APP
            )
        }

        if (newIcon == "default" || !iconAliases.contains(newIcon)) {
            val mainActivityName = ComponentName(packageName, "$packageName.MainActivity")
            packageManager.setComponentEnabledSetting(
                mainActivityName,
                PackageManager.COMPONENT_ENABLED_STATE_ENABLED,
                PackageManager.DONT_KILL_APP
            )
        } else {
            val newComponentName = ComponentName(packageName, "$packageName.$newIcon")
            packageManager.setComponentEnabledSetting(
                newComponentName,
                PackageManager.COMPONENT_ENABLED_STATE_ENABLED,
                PackageManager.DONT_KILL_APP
            )
        }
    }
}
