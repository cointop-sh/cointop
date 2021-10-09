// +build android

#include <android/log.h>
#include <jni.h>
#include <stdlib.h>
#include <string.h>

#define LOG_FATAL(...) __android_log_print(ANDROID_LOG_FATAL, "GoLog", __VA_ARGS__)

static const char *jstringToCharCopy(JNIEnv *env, const jstring string)
{
    const char *chars = (*env)->GetStringUTFChars(env, string, NULL);
    const char *copy = strdup(chars);
    (*env)->ReleaseStringUTFChars(env, string, chars);

    return copy;
}

static jclass findClass(JNIEnv *env, const char *class_name)
{
    jclass clazz = (*env)->FindClass(env, class_name);

    if (clazz == NULL)
    {
        (*env)->ExceptionClear(env);
        LOG_FATAL("cannot find %s", class_name);
        return NULL;
    }

    return clazz;
}

static jmethodID findMethod(JNIEnv *env, jclass clazz, const char *name, const char *sig)
{
    jmethodID m = (*env)->GetMethodID(env, clazz, name, sig);

    if (m == 0)
    {
        (*env)->ExceptionClear(env);
        LOG_FATAL("cannot find method %s %s", name, sig);
        return 0;
    }

    return m;
}

static jfieldID findField(JNIEnv *env, jclass clazz, const char *name, const char *sig)
{
    jfieldID f = (*env)->GetFieldID(env, clazz, name, sig);

    if (f == 0)
    {
        (*env)->ExceptionClear(env);
        LOG_FATAL("cannot find method %s %s", name, sig);
        return 0;
    }

    return f;
}

static jfieldID getStaticFieldID(JNIEnv *env, jclass clazz, const char *name, const char *sig)
{
    jfieldID f = (*env)->GetStaticFieldID(env, clazz, name, sig);

    if (f == 0)
    {
        (*env)->ExceptionClear(env);
        LOG_FATAL("cannot find static field %s %s", name, sig);
        return 0;
    }

    return f;
}

static const char *toLanguageTag(JNIEnv *env, jobject locale)
{
    const jclass java_util_Locale = findClass(env, "java/util/Locale");

    const jstring localeStr =
        (*env)->CallObjectMethod(
            env,
            locale,
            (*env)->GetMethodID(env, java_util_Locale, "toLanguageTag", "()Ljava/lang/String;"));

    return jstringToCharCopy(env, localeStr);
}

static const char *toLanguageTags(JNIEnv *env, jobject locales, jclass android_os_LocaleList)
{
    const jstring localeStr =
        (*env)->CallObjectMethod(
            env,
            locales,
            (*env)->GetMethodID(env, android_os_LocaleList, "toLanguageTags", "()Ljava/lang/String;"));

    return jstringToCharCopy(env, localeStr);
}

static int getAPIVersion(JNIEnv *env)
{
    // VERSION is a nested class within android.os.Build (hence "$" rather than "/")
    const jclass versionClass = findClass(env, "android/os/Build$VERSION");
    const jfieldID sdkIntFieldID = getStaticFieldID(env, versionClass, "SDK_INT", "I");

    int sdkInt = (*env)->GetStaticIntField(env, versionClass, sdkIntFieldID);

    return sdkInt;
}

static const jobject getConfiguration(JNIEnv *env, jobject context)
{
    const jclass android_content_ContextWrapper = findClass(env, "android/content/ContextWrapper");
    const jclass android_content_res_Resources = findClass(env, "android/content/res/Resources");

    const jobject resources =
        (*env)->CallObjectMethod(
            env,
            context,
            findMethod(env, android_content_ContextWrapper, "getResources", "()Landroid/content/res/Resources;"));
    const jobject configuration =
        (*env)->CallObjectMethod(
            env,
            resources,
            findMethod(env, android_content_res_Resources, "getConfiguration", "()Landroid/content/res/Configuration;"));

    return configuration;
}

static const jobject getLocaleObject(JNIEnv *env, jobject context)
{
    const jobject configuration = getConfiguration(env, context);
    const jclass android_content_res_Configuration = findClass(env, "android/content/res/Configuration");

    int version = getAPIVersion(env);

    // Android N or later
    // See https://developer.android.com/reference/android/content/res/Configuration#locale
    if (version >= 24) {
        const jclass android_os_LocaleList = findClass(env, "android/os/LocaleList");

        const jobject locales =
            (*env)->CallObjectMethod(
                env,
                configuration,
                findMethod(env, android_content_res_Configuration, "getLocales", "()Landroid/os/LocaleList;"));

        return (*env)->CallObjectMethod(
            env,
            locales,
            findMethod(env, android_os_LocaleList, "get", "(I)Ljava/util/Locale;"),
            0);
    } else {
        return (*env)->GetObjectField(
            env,
            configuration,
            findField(env, android_content_res_Configuration, "locale", "Ljava/util/Locale;"));
    }
}

// Basically the same as `getResources().getConfiguration().getLocales()` for Android N and later,
// or `getResources().getConfiguration().locale` for earlier Android version.
const char *getLocales(uintptr_t java_vm, uintptr_t jni_env, uintptr_t ctx)
{
    JavaVM *vm = (JavaVM *)java_vm;
    JNIEnv *env = (JNIEnv *)jni_env;
    jobject context = (jobject)ctx;

    const jobject configuration = getConfiguration(env, context);
    const jclass android_content_res_Configuration = findClass(env, "android/content/res/Configuration");

    int version = getAPIVersion(env);

    // Android N or later
    // See https://developer.android.com/reference/android/content/res/Configuration#locale
    if (version >= 24) {
        const jclass android_os_LocaleList = findClass(env, "android/os/LocaleList");

        const jobject locales =
            (*env)->CallObjectMethod(
                env,
                configuration,
                findMethod(env, android_content_res_Configuration, "getLocales", "()Landroid/os/LocaleList;"));

        return toLanguageTags(env, locales, android_os_LocaleList);
    } else {
        const jobject locale =
            (*env)->GetObjectField(
                env,
                configuration,
                findField(env, android_content_res_Configuration, "locale", "Ljava/util/Locale;"));

        return toLanguageTag(env, locale);
    }
}

// Basically the same as `getResources().getConfiguration().getLocales().get(0).toString()` for Android N and later,
// or `getResources().getConfiguration().locale` for earlier Android version.
const char *getLocale(uintptr_t java_vm, uintptr_t jni_env, uintptr_t ctx)
{
    JavaVM *vm = (JavaVM *)java_vm;
    JNIEnv *env = (JNIEnv *)jni_env;
    jobject context = (jobject)ctx;

    const jobject locale = getLocaleObject(env, context);

    return toLanguageTag(env, locale);
}

const char *getLanguage(uintptr_t java_vm, uintptr_t jni_env, uintptr_t ctx)
{
    JavaVM *vm = (JavaVM *)java_vm;
    JNIEnv *env = (JNIEnv *)jni_env;
    jobject context = (jobject)ctx;

    const jobject locale = getLocaleObject(env, context);
    const jclass java_util_Locale = findClass(env, "java/util/Locale");

    const jstring language =
        (*env)->CallObjectMethod(
            env,
            locale,
            (*env)->GetMethodID(env, java_util_Locale, "getLanguage", "()Ljava/lang/String;"));

    return jstringToCharCopy(env, language);
}

const char *getRegion(uintptr_t java_vm, uintptr_t jni_env, uintptr_t ctx)
{
    JavaVM *vm = (JavaVM *)java_vm;
    JNIEnv *env = (JNIEnv *)jni_env;
    jobject context = (jobject)ctx;

    const jobject locale = getLocaleObject(env, context);
    const jclass java_util_Locale = findClass(env, "java/util/Locale");

    const jstring country =
        (*env)->CallObjectMethod(
            env,
            locale,
            (*env)->GetMethodID(env, java_util_Locale, "getCountry", "()Ljava/lang/String;"));

    return jstringToCharCopy(env, country);
}
