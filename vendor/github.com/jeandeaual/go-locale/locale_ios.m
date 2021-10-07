// +build ios

#import <UIKit/UIKit.h>

const char *getLocale()
{
    NSString *locale = [[NSLocale preferredLanguages] firstObject];

    return [locale UTF8String];
}

const char *getLocales()
{
    NSString *locales = [[NSLocale preferredLanguages] componentsJoinedByString:@","];

    return [locales UTF8String];
}
