int violationPolicyNotify(void *p)
{
    int ViolationPolicyRegistration(void *);
    return ViolationPolicyRegistration(p);
}

int voidPolicyCallback(void *p)
{
    int VoidPolicyCallback(void *);
    return VoidPolicyCallback(p);
}