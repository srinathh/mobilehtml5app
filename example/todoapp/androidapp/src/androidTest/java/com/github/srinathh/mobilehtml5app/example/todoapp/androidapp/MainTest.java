package com.github.srinathh.mobilehtml5app.example.todoapp.androidapp;

import android.test.ActivityInstrumentationTestCase2;

/**
 * This is a simple framework for a test of an Application.  See
 * {@link android.test.ApplicationTestCase ApplicationTestCase} for more information on
 * how to write and extend Application tests.
 * <p/>
 * To run this test, you can type:
 * adb shell am instrument -w \
 * -e class com.github.srinathh.mobilehtml5app.example.todoapp.androidapp.MainTest \
 * com.github.srinathh.mobilehtml5app.example.todoapp.androidapp.tests/android.test.InstrumentationTestRunner
 */
public class MainTest extends ActivityInstrumentationTestCase2<Main> {

    public MainTest() {
        super("com.github.srinathh.mobilehtml5app.example.todoapp.androidapp", Main.class);
    }

}
