// adapted from https://github.com/slymax/webview/blob/master/app/src/main/java/com/example/app/MyAppWebViewClient.java
// TODO: Same license, no name in copyright notice on source; make sure I've attributed properly

package com.midnightfreddie.editformcpe;

import android.content.Intent;
import android.net.Uri;
import android.webkit.WebView;
import android.webkit.WebViewClient;

public class MyAppWebViewClient extends WebViewClient {

    @Override
    public boolean shouldOverrideUrlLoading(WebView view, String url) {
        if (Uri.parse(url).getHost().endsWith("example.com")) {
            return false;
        }

        Intent intent = new Intent(Intent.ACTION_VIEW, Uri.parse(url));
        view.getContext().startActivity(intent);
        return true;
    }
}